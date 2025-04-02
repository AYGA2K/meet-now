<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <!-- Room Controls -->
    <div class="bg-gray-800 p-4">
      <div class="flex items-center mb-4">
        <input v-model="userId" placeholder="Enter User ID" class="mr-2 p-2 bg-gray-700 rounded">
        <input v-model="roomId" placeholder="Enter room ID" class="mr-2 p-2 bg-gray-700 rounded">
        <button class="bg-blue-600 px-4 py-2 rounded" @click="joinRoom">
          Join Room
        </button>
      </div>

      <div class="grid grid-cols-2 gap-4 mb-4">
        <div>
          <h2 class="text-xl mb-2">Your Media</h2>
          <video v-if="hasVideo" id="localVideo" autoplay muted playsinline class="w-full bg-black rounded" />
          <audio v-else id="localAudio" autoplay muted class="w-full bg-black rounded" />
          <div class="mt-2 flex space-x-2">
            <button class="bg-gray-700 px-3 py-1 rounded" @click="toggleMute">
              {{ isMuted ? "Unmute" : "Mute" }}
            </button>
            <button v-if="hasVideo" class="bg-gray-700 px-3 py-1 rounded" @click="toggleVideo">
              {{ isVideoOff ? "Show Video" : "Hide Video" }}
            </button>
          </div>
        </div>
        <div>
          <h2 class="text-xl mb-2">Remote Media</h2>
          <video v-if="hasVideo" id="remoteVideo" autoplay playsinline class="w-full bg-black rounded" />
          <audio v-else id="remoteAudio" autoplay class="w-full bg-black rounded" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";

// WebSocket State
const socket = ref<WebSocket | null>(null);
const userId = ref<string>(""); // Replace with a unique user ID
const roomId = ref<string>("");
const peerConnection = ref<RTCPeerConnection | null>(null);
const localStream = ref<MediaStream | null>(null);

// UI State
const isMuted = ref(true);
const isVideoOff = ref(false);
const hasVideo = ref(true); // Set to false if you want audio-only mode

const iceServers = [
  { urls: "stun:stun.l.google.com:19302" },
  { urls: "stun:stun1.l.google.com:19302" },
];

// Initialize media stream
async function initMedia() {
  try {
    localStream.value = await navigator.mediaDevices.getUserMedia({
      audio: true,
      video: hasVideo.value,
    });

    const localVideo = document.getElementById(
      "localVideo",
    ) as HTMLVideoElement;
    const localAudio = document.getElementById(
      "localAudio",
    ) as HTMLAudioElement;

    if (hasVideo.value && localVideo) {
      localVideo.srcObject = localStream.value;
    } else if (!hasVideo.value && localAudio) {
      localAudio.srcObject = localStream.value;
    }

    toggleMute();
  } catch (error) {
    console.error("Error accessing media devices:", error);
    alert("Failed to access media devices. Please check permissions.");
  }
}

// Toggle audio mute
function toggleMute() {
  if (localStream.value) {
    localStream.value.getAudioTracks().forEach((track) => {
      track.enabled = !track.enabled;
    });
    isMuted.value = !isMuted.value;
  }
}

// Toggle video on/off
function toggleVideo() {
  if (localStream.value) {
    localStream.value.getVideoTracks().forEach((track) => {
      track.enabled = !track.enabled;
    });
    isVideoOff.value = !isVideoOff.value;
  }
}

// Create peer connection
function createPeerConnection() {
  peerConnection.value = new RTCPeerConnection({ iceServers });

  if (localStream.value) {
    localStream.value.getTracks().forEach((track) => {
      peerConnection.value?.addTrack(track, localStream.value!);
    });
  }

  peerConnection.value.ontrack = (event) => {
    const remoteVideo = document.getElementById(
      "remoteVideo",
    ) as HTMLVideoElement;
    const remoteAudio = document.getElementById(
      "remoteAudio",
    ) as HTMLAudioElement;

    if (hasVideo.value && remoteVideo && event.streams[0]) {
      remoteVideo.srcObject = event.streams[0];
    } else if (!hasVideo.value && remoteAudio && event.streams[0]) {
      remoteAudio.srcObject = event.streams[0];
    }
  };

  peerConnection.value.onicecandidate = (event) => {
    if (event.candidate) {
      sendMessage({
        type: "candidate",
        candidate: JSON.stringify(event.candidate),
        from: userId.value,
        room_id: roomId.value,
      });
    }
  };
}

// Send message via WebSocket
function sendMessage(message: unknown) {
  if (socket.value && socket.value.readyState === WebSocket.OPEN) {
    socket.value.send(JSON.stringify(message));
  }
}

// Join room
function joinRoom() {
  if (!roomId.value) {
    alert("Please enter a room ID.");
    return;
  }

  connect();

  // Notify the server that the user has joined the room
  sendMessage({
    type: "join-room",
    from: userId.value,
    room_id: roomId.value,
  });

  // Start the WebRTC connection
  createPeerConnection();
  startConnection();
}

// Start connection as caller
async function startConnection() {
  try {
    const offer = await peerConnection.value!.createOffer();
    await peerConnection.value!.setLocalDescription(offer);

    sendMessage({
      type: "offer",
      sdp: JSON.stringify(peerConnection.value!.localDescription),
      from: userId.value,
      room_id: roomId.value,
    });
  } catch (error) {
    console.error("Error creating offer:", error);
    alert("Failed to start the connection.");
  }
}

// Handle incoming SDP offer
async function handleOffer(message: unknown) {
  if (!peerConnection.value) {
    createPeerConnection();
  }

  const offer = new RTCSessionDescription(JSON.parse(message.sdp));
  await peerConnection.value.setRemoteDescription(offer);

  const answer = await peerConnection.value.createAnswer();
  await peerConnection.value.setLocalDescription(answer);

  sendMessage({
    type: "answer",
    sdp: JSON.stringify(peerConnection.value.localDescription),
    from: userId.value,
    room_id: roomId.value,
  });
}

// Handle incoming SDP answer
async function handleAnswer(message: any) {
  const answer = new RTCSessionDescription(JSON.parse(message.sdp));
  await peerConnection.value?.setRemoteDescription(answer);
}

// Handle incoming ICE candidate
async function handleCandidate(message: any) {
  if (peerConnection.value) {
    const candidate = new RTCIceCandidate(JSON.parse(message.candidate));
    await peerConnection.value.addIceCandidate(candidate);
  }
}

// WebSocket Connection
function connect(): void {
  if (socket.value) {
    console.warn("Already connected!");
    return;
  }

  socket.value = new WebSocket(`ws://localhost:42157/ws/${userId.value}`);

  socket.value.onmessage = (event: MessageEvent) => {
    const message = JSON.parse(event.data);

    switch (message.type) {
      case "offer":
        handleOffer(message);
        break;
      case "answer":
        handleAnswer(message);
        break;
      case "candidate":
        handleCandidate(message);
        break;
      default:
        console.warn("Unknown message type:", message.type);
    }
  };

  socket.value.onopen = () => {
    console.log("Connected to WebSocket server");
  };

  socket.value.onclose = () => {
    console.log("Disconnected from WebSocket server");
    socket.value = null;
  };
}

// Lifecycle Hooks
onMounted(() => {
  initMedia();
});

onUnmounted(() => {
  if (socket.value) {
    socket.value.close();
  }

  if (peerConnection.value) {
    peerConnection.value.close();
  }
});
</script>
