import {defineStore} from "pinia";
import {ref, computed} from "vue";

export interface Hunk {
  header: string;
  lines: string[];
}

export const useReviewStore = defineStore("review", () => {
  const gitDiff = ref("");
  const parsedDiff = computed<Hunk[]>(() => {
    if (!gitDiff.value) return [];
    // Parsing logic here...
    return [];
  });

  function setGeneratedDiff(diff: string) {
    gitDiff.value = diff;
  }

  return {gitDiff, parsedDiff, setGeneratedDiff};
});