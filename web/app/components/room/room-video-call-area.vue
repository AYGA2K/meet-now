<template>
  <div
    class="flex-1 aspect-video justify-center items-center p-2 flex flex-col overflow-hidden"
  >
    <!-- Main Video Display -->
    <RoomMainVideo :is-muted="isMuted" :is-screen-sharing="isScreenSharing" />

    <!-- Participants Grid -->
    <RoomParticipants :participants="participants" />

    <!-- Call Controls -->
    <RoomCallControll
      :is-muted="isMuted"
      :is-video-off="isVideoOff"
      :is-screen-sharing="isScreenSharing"
      :is-chat-open="isChatOpen"
      :is-participants-open="isParticipantsOpen"
      @toggle-mute="$emit('toggle-mute')"
      @toggle-video="$emit('toggle-video')"
      @toggle-screen-share="$emit('toggle-screen-share')"
      @toggle-chat="$emit('toggle-chat')"
      @toggle-participants="$emit('toggle-participants')"
    />
  </div>
</template>

<script setup lang="ts">
import { RoomCallControll, RoomMainVideo, RoomParticipants } from "#components";
import type { Participant } from "~/types";

defineProps({
  isMuted: Boolean,
  isVideoOff: Boolean,
  isScreenSharing: Boolean,
  isChatOpen: Boolean,
  isParticipantsOpen: Boolean,
  participants: {
    type: Array as () => Participant[],
    default: () => [],
  },
});

defineEmits([
  "toggle-mute",
  "toggle-video",
  "toggle-screen-share",
  "toggle-chat",
  "toggle-participants",
]);
</script>
