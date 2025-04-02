<template>
  <div class="grid grid-cols-2 md:grid-cols-5 gap-1 mb-4 h-40">
    <!-- Dynamic Participants -->
    <div
      v-for="(participant, index) in participants"
      :key="index"
      class="aspect-video relative bg-gray-800 rounded-lg overflow-hidden h-full"
    >
      <img
        :src="`https://i.pravatar.cc/300?img=${participant.avatar}`"
        :alt="participant.name"
        class="w-full h-full object-cover"
      />
      <div
        class="absolute bottom-2 left-2 bg-gray-800 bg-opacity-70 px-2 py-0.5 rounded text-white text-xs flex items-center"
      >
        <UIcon
          v-if="participant.isMuted"
          name="i-lucide-mic-off"
          class="mr-1 text-xs"
        />
        <span>{{ participant.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Participant } from "~/types";

defineProps({
  participants: {
    type: Array as () => Participant[],
    required: true,
  },
});
</script>
