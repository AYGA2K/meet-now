<template>
  <ul class="space-y-3">
    <li
      v-for="(participant, index) in participants"
      :key="index"
      class="flex items-center justify-between"
    >
      <div class="flex items-center">
        <UAvatar
          :src="`https://i.pravatar.cc/150?img=${participant.avatar}`"
          :alt="participant.name"
          size="sm"
          class="mr-2"
        />
        <span class="text-sm">{{ participant.name }}</span>
      </div>
      <div class="flex items-center space-x-1">
        <UIcon v-if="participant.isMuted" name="i-lucide-mic-off" />
        <UIcon
          v-if="participant.isHost"
          name="i-lucide-shield"
          class="text-emerald-500"
        />
      </div>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import type { Participant } from "~/types";

defineProps({
  participants: {
    type: Array as () => Participant[],
    required: true,
  },
});
</script>
