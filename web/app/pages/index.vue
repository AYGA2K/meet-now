<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <!-- Connection Controls -->
    <div class="bg-gray-800 p-4">
      <div class="flex items-center mb-4">
        <input id="isCaller" v-model="isCaller" type="checkbox" class="mr-2" />
        <label for="isCaller">I am the caller</label>
      </div>
      <div class="grid grid-cols-2 gap-4 mb-4">
        <div>
          <label class="block mb-2">Local SDP</label>
          <textarea
            v-model="localSDP"
            class="w-full bg-gray-700 p-2 rounded"
            rows="4"
            readonly
          />
          <button
            class="mt-2 bg-blue-600 px-4 py-2 rounded"
            @click="copyLocalSDP"
          >
            Copy SDP
          </button>
        </div>
        <div>
          <label class="block mb-2">Remote SDP</label>
          <textarea
            v-model="remoteSDP"
            class="w-full bg-gray-700 p-2 rounded"
            rows="4"
          />
          <button
            class="mt-2 bg-green-600 px-4 py-2 rounded"
            @click="setRemoteSDP"
          >
            Set Remote SDP
          </button>
        </div>
      </div>

      <div class="flex space-x-4">
        <button
          class="bg-purple-600 px-4 py-2 rounded"
          :disabled="!isCaller"
          @click="startConnection"
        >
          Start Call
        </button>
        <button
          class="bg-orange-600 px-4 py-2 rounded"
          :disabled="isCaller"
          @click="acceptConnection"
        >
          Accept Call
        </button>
      </div>
    </div>

    <!-- Video Areas -->
    <div class="grid grid-cols-2 gap-4 p-4">
      <div>
        <h2 class="text-xl mb-2">Your Video</h2>
        <video
          id="localVideo"
          autoplay
          playsinline
          muted
          class="w-full bg-black rounded"
        />
        <div class="mt-2">
          <button
            class="mr-2 bg-gray-700 px-3 py-1 rounded"
            @click="toggleMute"
          >
            {{ isMuted ? "Unmute" : "Mute" }}
          </button>
          <button class="bg-gray-700 px-3 py-1 rounded" @click="toggleVideo">
            {{ isVideoOff ? "Show Video" : "Hide Video" }}
          </button>
        </div>
      </div>
      <div>
        <h2 class="text-xl mb-2">Remote Video</h2>
        <video
          id="remoteVideo"
          autoplay
          playsinline
          class="w-full bg-black rounded"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";

// UI State
const isMuted = ref(true);
const isVideoOff = ref(false);
const isCaller = ref(true);
const localSDP = ref("");
const remoteSDP = ref("");

// WebRTC Variables
let localStream: MediaStream | null = null;
let peerConnection: RTCPeerConnection | null = null;

const iceServers = [
  { urls: "stun:stun.l.google.com:19302" },
  { urls: "stun:stun1.l.google.com:19302" },
];

// Initialize media stream
async function initMedia() {
  try {
    localStream = await navigator.mediaDevices.getUserMedia({
      audio: true,
      video: true,
    });

    const localVideo = document.getElementById(
      "localVideo",
    ) as HTMLVideoElement;
    if (localVideo) {
      localVideo.srcObject = localStream;
    }

    // Start muted by default
    toggleMute();
  } catch (error) {
    console.error("Error accessing media devices:", error);
  }
}
// Toggle audio mute
function toggleMute() {
  if (localStream) {
    localStream.getAudioTracks().forEach((track) => {
      track.enabled = !track.enabled;
    });
    isMuted.value = !isMuted.value;
  }
}

// Toggle video
function toggleVideo() {
  if (localStream) {
    localStream.getVideoTracks().forEach((track) => {
      track.enabled = !track.enabled;
    });
    isVideoOff.value = !isVideoOff.value;
  }
}

// Create peer connection
function createPeerConnection() {
  peerConnection = new RTCPeerConnection({ iceServers });

  // Add local stream to connection
  if (localStream) {
    localStream.getTracks().forEach((track) => {
      peerConnection!.addTrack(track, localStream!);
    });
  }

  // Handle remote stream
  peerConnection.ontrack = (event) => {
    const remoteVideo = document.getElementById(
      "remoteVideo",
    ) as HTMLVideoElement;
    if (remoteVideo) {
      if (event.streams[0] != null) {
        remoteVideo.srcObject = event.streams[0];
      }
    }
  };

  // ICE candidate handling
  peerConnection.onicecandidate = (event) => {
    if (event.candidate) {
      localSDP.value = JSON.stringify(peerConnection!.localDescription);
    }
  };
}

// Start connection as caller
async function startConnection() {
  createPeerConnection();
  try {
    const offer = await peerConnection!.createOffer();
    await peerConnection!.setLocalDescription(offer);
    localSDP.value = JSON.stringify(offer);
  } catch (error) {
    console.error("Error creating offer:", error);
  }
}

// Accept connection as callee
async function acceptConnection() {
  if (!remoteSDP.value) return;
  createPeerConnection();
  try {
    const remoteDesc = JSON.parse(remoteSDP.value);
    await peerConnection!.setRemoteDescription(remoteDesc);
    const answer = await peerConnection!.createAnswer();
    await peerConnection!.setLocalDescription(answer);
    localSDP.value = JSON.stringify(answer);
  } catch (error) {
    console.error("Error accepting connection:", error);
  }
}

// Set remote SDP
async function setRemoteSDP() {
  if (!peerConnection || !remoteSDP.value) return;
  try {
    const remoteDesc = JSON.parse(remoteSDP.value);
    await peerConnection.setRemoteDescription(remoteDesc);
  } catch (error) {
    console.error("Error setting remote description:", error);
  }
}

// Copy local SDP to clipboard
function copyLocalSDP() {
  navigator.clipboard.writeText(localSDP.value);
  alert("SDP copied to clipboard!");
}

onMounted(() => {
  initMedia();
  $nuxtSocket
});
</script>
