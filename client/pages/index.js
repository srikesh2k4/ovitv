import { useEffect, useRef, useState } from 'react';

export default function Home(){
  const localRef = useRef();
  const remoteRef = useRef();
  const pcRef = useRef(null);
  const wsRef = useRef(null);
  const [status, setStatus] = useState('init');
  const [muted, setMuted] = useState(false);
  const [videoOn, setVideoOn] = useState(true);

  useEffect(() => {
    // request media immediately
    (async ()=>{
      try{
        const s = await navigator.mediaDevices.getUserMedia({ video:true, audio:true });
        localRef.current.srcObject = s;
        localRef.current.muted = true;
        setStatus('ready');
        pcRef.current = createPeerConnection();
      }catch(err){
        console.error(err);
        setStatus('no-media');
      }
    })();

    // cleanup
    return ()=> {
      if (pcRef.current) { pcRef.current.close(); pcRef.current = null; }
      if (wsRef.current) { wsRef.current.close(); wsRef.current = null; }
    }
  }, []);

  function createPeerConnection(){
    const pc = new RTCPeerConnection();
    pc.ontrack = (e) => {
      if (remoteRef.current) remoteRef.current.srcObject = e.streams[0];
    };
    pc.onicecandidate = (e) => {
      if (e.candidate && wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'ice', candidate: e.candidate }));
      }
    };
    // add local tracks
    if (localRef.current && localRef.current.srcObject) {
      localRef.current.srcObject.getTracks().forEach(t => pc.addTrack(t, localRef.current.srcObject));
    }
    return pc;
  }

  async function connectSignaling(){
    setStatus('connecting');
    const protocol = location.protocol === 'https:' ? 'wss' : 'ws';
    const url = `${protocol}://${location.host}/signaling`;
    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => setStatus('connected');
    ws.onmessage = async (evt) => {
      const data = JSON.parse(evt.data);
      console.log('ws recv', data);
      if (data.type === 'paired') {
        setStatus('paired');
        // if this client should create offer -> create and send
        // We'll decide: first one to receive 'paired' will create an offer after short delay to avoid collision.
        setTimeout(()=>maybeCreateOffer(), 200);
      } else if (data.type === 'waiting') {
        setStatus('waiting');
      } else if (data.type === 'offer') {
        await handleOffer(data.offer);
      } else if (data.type === 'answer') {
        await pcRef.current.setRemoteDescription(new RTCSessionDescription(data.answer));
      } else if (data.type === 'ice') {
        if (data.candidate) {
          try { await pcRef.current.addIceCandidate(data.candidate); } catch(e){ console.warn(e); }
        }
      } else if (data.type === 'partner-left') {
        setStatus('partner-left');
        // remote video cleared
        if (remoteRef.current) remoteRef.current.srcObject = null;
      } else if (data.type === 'no-peer') {
        setStatus('no-peer');
      }
    };

    ws.onclose = ()=> setStatus('closed');
    ws.onerror = (e)=> console.error(e);
  }

  async function maybeCreateOffer(){
    try{
      if (!pcRef.current) pcRef.current = createPeerConnection();
      const offer = await pcRef.current.createOffer();
      await pcRef.current.setLocalDescription(offer);
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'offer', offer }));
      }
    }catch(e){ console.error(e); }
  }

  async function handleOffer(offer){
    try{
      if (!pcRef.current) pcRef.current = createPeerConnection();
      await pcRef.current.setRemoteDescription(new RTCSessionDescription(offer));
      const answer = await pcRef.current.createAnswer();
      await pcRef.current.setLocalDescription(answer);
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'answer', answer }));
      }
    }catch(e){ console.error(e); }
  }

  function toggleMute(){
    if (!localRef.current || !localRef.current.srcObject) return;
    localRef.current.srcObject.getAudioTracks().forEach(t => t.enabled = !t.enabled);
    setMuted(prev=>!prev);
  }
  function toggleVideo(){
    if (!localRef.current || !localRef.current.srcObject) return;
    localRef.current.srcObject.getVideoTracks().forEach(t => t.enabled = !t.enabled);
    setVideoOn(prev=>!prev);
  }
  function nextPartner(){
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type: 'next' }));
      // clear remote
      if (remoteRef.current) remoteRef.current.srcObject = null;
      setStatus('looking');
    }
  }

  return (
    <div className="container">
      <div className="header">
        <h1>Mini OmeTV — deployable on Railway</h1>
        <div className="controls">
          <button className="btn" onClick={connectSignaling}>Connect</button>
          <button className="btn ghost" onClick={toggleMute}>{muted ? 'Unmute' : 'Mute'}</button>
          <button className="btn ghost" onClick={toggleVideo}>{videoOn ? 'Video Off' : 'Video On'}</button>
          <button className="btn" onClick={nextPartner}>Next</button>
        </div>
      </div>

      <div className="status">Status: {status}</div>

      <div className="video-grid" style={{marginTop:20}}>
        <div>
          <h3 style={{margin:6}}>You</h3>
          <video ref={localRef} autoPlay playsInline></video>
        </div>
        <div>
          <h3 style={{margin:6}}>Partner</h3>
          <video ref={remoteRef} autoPlay playsInline></video>
        </div>
      </div>

      <div style={{marginTop:20}}>
        <p>Report abusive user: <small>use /api/report with {`{ reason, details }`}</small></p>
      </div>
    </div>
  )
}
