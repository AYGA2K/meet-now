<template>
  <div class="bg-gray-900 flex flex-col overflow-hidden">
    <!-- Header -->
    <RoomHeader />

    <!-- Manual Connection Controls (for testing) -->
    <div class="bg-gray-800 p-4 text-white">
      <div class="flex items-center mb-2">
        <input id="isCaller" v-model="isCaller" type="checkbox" class="mr-2" />
        <label for="isCaller">I am the caller</label>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block mb-1">Offer SDP</label>
          <textarea
            v-model="offerInput"
            class="w-full bg-gray-700 text-white p-2 rounded mb-2"
            rows="4"
          />
          <button
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded"
            :disabled="isCaller"
            @click="setRemoteOffer"
          >
            Set Remote Offer
          </button>
        </div>

        <div>
          <label class="block mb-1">Answer SDP</label>
          <textarea
            v-model="answerInput"
            class="w-full bg-gray-700 text-white p-2 rounded mb-2"
            rows="4"
          />
          <button
            class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded"
            :disabled="!isCaller"
            @click="setRemoteAnswer"
          >
            Set Remote Answer
          </button>
        </div>
      </div>

      <div class="mt-4 flex space-x-4">
        <button
          class="bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded"
          :disabled="!isCaller"
          @click="startCall"
        >
          Start Call
        </button>
        <button
          class="bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded"
          :disabled="isCaller || !offerInput"
          @click="answerCall"
        >
          Answer Call
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <main class="flex-1 flex overflow-hidden">
      <!-- Video Call Area -->
      <RoomVideoCallArea
        :is-muted="isMuted"
        :is-video-off="isVideoOff"
        :is-screen-sharing="isScreenSharing"
        :is-chat-open="isChatOpen"
        :is-participants-open="isParticipantsOpen"
        :participants="participants"
        @toggle-mute="toggleMute"
        @toggle-video="toggleVideo"
        @toggle-screen-share="toggleScreenShare"
        @toggle-chat="isChatOpen = !isChatOpen"
        @toggle-participants="isParticipantsOpen = !isParticipantsOpen"
      />

      <!-- Sidebar (Chat or Participants) -->
      <RoomSidebar
        v-if="isChatOpen || isParticipantsOpen"
        :is-chat-open="isChatOpen"
        :is-participants-open="isParticipantsOpen"
        :chat-messages="chatMessages"
        :participants="participants"
        @close-sidebar="
          isChatOpen = false;
          isParticipantsOpen = false;
        "
      />
    </main>
  </div>
</template>

<script setup lang="ts">
import { RoomHeader, RoomSidebar, RoomVideoCallArea } from "#components";
import type { ChatMessage, Participant, RTCConfiguration } from "~/types";

// UI State
const isMuted = ref(true);
const isVideoOff = ref(true);
const isScreenSharing = ref(false);
const isChatOpen = ref(false);
const isParticipantsOpen = ref(false);
const isCaller = ref(true);
const offerInput = ref("");
const answerInput = ref("");

// Sample Data
const chatMessages = ref<ChatMessage[]>([
  {
    sender: "John Doe",
    text: "Hello everyone!",
    time: "10:30 AM",
    avatar: "32",
  },
  {
    sender: "Sarah Johnson",
    text: "Hi John, how are you?",
    time: "10:31 AM",
    avatar: "44",
  },
]);

const participants = ref<Participant[]>([
  { name: "You (Host)", avatar: "68", isMuted: false, isHost: true },
  { name: "John Doe", avatar: "32", isMuted: false, isHost: false },
]);

// WebRTC Configuration
const servers = {
  iceServers: [
    { urls: "stun:stun.l.google.com:19302" },
    { urls: "stun:stun1.l.google.com:19302" },
  ],
};

let localStream: MediaStream | null = null;
let screenStream: MediaStream | null = null;
let remoteStream: MediaStream | null = null;
let localPC: RTCPeerConnection | null = null;
let remotePC: RTCPeerConnection | null = null;

async function getMediaStream() {
  try {
    localStream = await navigator.mediaDevices.getUserMedia({
      audio: true,
      video: true,
    });

    const videoElement = document.getElementById(
      isCaller.value ? "you" : "other",
    ) as HTMLVideoElement;

    if (videoElement) {
      videoElement.srcObject = localStream;
    }

    toggleMute();
  } catch (error) {
    console.error("Error accessing media devices:", error);
  }
}

function toggleMute() {
  if (localStream) {
    localStream
      .getAudioTracks()
      .forEach((track) => (track.enabled = !isMuted.value));
    isMuted.value = !isMuted.value;
  }
}

function toggleVideo() {
  if (localStream) {
    localStream
      .getVideoTracks()
      .forEach((track) => (track.enabled = !isVideoOff.value));
    isVideoOff.value = !isVideoOff.value;
  }
}

async function toggleScreenShare() {
  if (!isScreenSharing.value) {
    try {
      screenStream = await navigator.mediaDevices.getDisplayMedia({
        video: true,
      });
      const videoElement = document.getElementById(
        isCaller.value ? "you" : "other",
      ) as HTMLVideoElement;
      if (videoElement) {
        videoElement.srcObject = screenStream;
      }
      isScreenSharing.value = true;
    } catch (error) {
      console.error("Error accessing screen share:", error);
    }
  } else {
    if (screenStream) {
      screenStream.getTracks().forEach((track) => track.stop());
    }
    if (localStream) {
      const videoElement = document.getElementById(
        isCaller.value ? "you" : "other",
      ) as HTMLVideoElement;
      if (videoElement) {
        videoElement.srcObject = localStream;
      }
    }
    isScreenSharing.value = false;
  }
}

// WebRTC Functions
async function startCall() {
  if (!localStream) return;

  localPC = new RTCPeerConnection(servers);
  remotePC = new RTCPeerConnection(servers);

  // Add local stream to connection
  localStream.getTracks().forEach((track) => {
    localPC!.addTrack(track, localStream!);
  });

  // Set up remote connection
  remotePC.ontrack = (event) => {
    const videoElement = document.getElementById(
      isCaller.value ? "other" : "you",
    ) as HTMLVideoElement;
    if (!remoteStream) {
      remoteStream = new MediaStream();
    }
    remoteStream.addTrack(event.track);
    videoElement.srcObject = remoteStream;
  };

  // ICE candidate exchange
  localPC.onicecandidate = (event) => {
    if (event.candidate) {
      remotePC!.addIceCandidate(event.candidate);
    }
  };

  remotePC.onicecandidate = (event) => {
    if (event.candidate) {
      localPC!.addIceCandidate(event.candidate);
    }
  };

  // Create offer
  try {
    const offer = await localPC.createOffer();
    await localPC.setLocalDescription(offer);

    // For manual testing, show the offer in the UI
    offerInput.value = JSON.stringify(localPC.localDescription);

    console.log(
      "Call started successfully - copy the Offer SDP to the other user",
    );
  } catch (error) {
    console.error("Error starting call:", error);
  }
}

async function answerCall() {
  if (!localStream || !offerInput.value) return;

  localPC = new RTCPeerConnection(servers);
  remotePC = new RTCPeerConnection(servers);

  // Add local stream to connection
  localStream.getTracks().forEach((track) => {
    localPC!.addTrack(track, localStream!);
  });

  // Set up remote connection
  remotePC.ontrack = (event) => {
    const videoElement = document.getElementById(
      isCaller.value ? "other" : "you",
    ) as HTMLVideoElement;
    if (!remoteStream) {
      remoteStream = new MediaStream();
    }
    remoteStream.addTrack(event.track);
    videoElement.srcObject = remoteStream;
  };

  // ICE candidate exchange
  localPC.onicecandidate = (event) => {
    if (event.candidate) {
      remotePC!.addIceCandidate(event.candidate);
    }
  };

  remotePC.onicecandidate = (event) => {
    if (event.candidate) {
      localPC!.addIceCandidate(event.candidate);
    }
  };

  try {
    // Set remote description from offer
    await localPC.setRemoteDescription(JSON.parse(offerInput.value));

    // Create answer
    const answerDescription = await localPC.createAnswer();
    await localPC.setLocalDescription(answerDescription);

    // Show the answer in the UI
    answerInput.value = JSON.stringify(answerDescription);

    console.log(
      "Call answered successfully - copy the Answer SDP back to the caller",
    );
  } catch (error) {
    console.error("Error answering call:", error);
  }
}

async function setRemoteOffer() {
  if (!localPC || !offerInput.value) return;
  try {
    await localPC.setRemoteDescription(JSON.parse(offerInput.value));
    console.log("Remote offer set successfully");
  } catch (error) {
    console.error("Error setting remote offer:", error);
  }
}

async function setRemoteAnswer() {
  if (!localPC || !answerInput.value) return;
  try {
    await localPC.setRemoteDescription(JSON.parse(answerInput.value));
    console.log("Remote answer set successfully");
  } catch (error) {
    console.error("Error setting remote answer:", error);
  }
}

onMounted(getMediaStream);
</script>

<style scoped>
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
