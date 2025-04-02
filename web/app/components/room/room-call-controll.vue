<template>
  <div class="rounded-lg p-4 shadow-sm mt-auto">
    <div class="flex justify-center items-center space-x-4">
      <UButton
        :icon="!isMuted ? 'i-lucide-mic-off' : 'i-lucide-mic'"
        :color="isMuted ? 'info' : 'error'"
        variant="soft"
        class="rounded-full h-12 w-12 flex items-center justify-center"
        aria-label="Toggle microphone"
        @click="$emit('toggle-mute')"
      />
      <UButton
        :icon="!isVideoOff ? 'i-lucide-video-off' : 'i-lucide-video'"
        :color="isVideoOff ? 'info' : 'error'"
        variant="soft"
        class="rounded-full h-12 w-12 flex items-center justify-center"
        aria-label="Toggle video"
        @click="$emit('toggle-video')"
      />
      <UButton
        icon="i-lucide-monitor"
        :color="isScreenSharing ? 'info' : 'error'"
        variant="soft"
        class="rounded-full h-12 w-12 flex items-center justify-center"
        aria-label="Share screen"
        @click="$emit('toggle-screen-share')"
      />
      <UButton
        icon="i-lucide-message-square"
        color="primary"
        variant="soft"
        class="rounded-full h-12 w-12 flex items-center justify-center"
        aria-label="Open chat"
        @click="$emit('toggle-chat')"
      />
      <UButton
        icon="i-lucide-users"
        color="primary"
        variant="soft"
        class="rounded-full h-12 w-12 flex items-center justify-center"
        aria-label="Show participants"
        @click="$emit('toggle-participants')"
      />
      <UButton
        icon="i-lucide-phone-off"
        color="error"
        variant="solid"
        class="rounded-full h-12 w-12 flex items-center justify-center"
        aria-label="End call"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  isMuted: Boolean,
  isVideoOff: Boolean,
  isScreenSharing: Boolean,
  isChatOpen: Boolean,
  isParticipantsOpen: Boolean,
});

defineEmits([
  "toggle-mute",
  "toggle-video",
  "toggle-screen-share",
  "toggle-chat",
  "toggle-participants",
]);
</script>
