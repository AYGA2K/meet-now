<template>
  <div class="relative flex-1 bg-gray-900 aspect-video rounded-lg overflow-hidden mb-4">
    <!-- <video id="video-main" ref="videoMain" autoplay /> -->
    <video id="other" autoplay />

    <div v-if="!isMuted"
      class="absolute bottom-4 left-4 bg-opacity-70 px-3 py-1 rounded-md text-white flex items-center">
      <UIcon name="i-lucide-mic-off" class="mr-2" />
      <span>Muted</span>
    </div>

    <div v-if="isScreenSharing"
      class="absolute top-4 left-4 bg-emerald-600 bg-opacity-90 px-3 py-1 rounded-md text-white flex items-center">
      <UIcon name="i-lucide-monitor" class="mr-2" />
      <span>You are sharing your screen</span>
    </div>

    <div class="absolute bottom-4 right-4 w-48 h-36 bg-gray-800 rounded-lg overflow-hidden border-2 border-white">
      <video id="you" autoplay />
      <div class="absolute bottom-2 left-2 bg-gray-800 bg-opacity-70 px-2 py-0.5 rounded text-white text-xs">
        You
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  isScreenSharing: Boolean,
  isMuted: Boolean,
  videoRef: {
    type: String,
    default: "video-local",
  },
});
</script>
