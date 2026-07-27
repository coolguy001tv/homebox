<template>
  <BaseModal v-model="open">
    <template #title> Browse NAS Files </template>

    <div class="flex flex-col gap-3 min-h-[300px] max-h-[60vh]">
      <!-- Breadcrumb navigation -->
      <div class="flex items-center gap-1 text-sm breadcrumbs overflow-x-auto">
        <ul class="flex items-center gap-0">
          <li>
            <button class="link link-hover" @click="goToDir(-1)">Home</button>
          </li>
          <li v-for="(crumb, idx) in breadcrumbs" :key="idx" class="flex items-center gap-0">
            <span class="mx-1">/</span>
            <button class="link link-hover truncate max-w-[150px]" @click="goToDir(idx)">
              {{ crumb }}
            </button>
          </li>
        </ul>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center items-center flex-1">
        <span class="loading loading-spinner loading-md"></span>
      </div>

      <!-- Error -->
      <div v-else-if="errorMsg" class="alert alert-error flex-1">
        <Icon name="mdi-alert-circle" class="h-5 w-5" />
        <span>{{ errorMsg }}</span>
      </div>

      <!-- Empty -->
      <div v-else-if="dirs.length === 0 && files.length === 0" class="flex justify-center items-center flex-1 text-gray-500">
        <div class="text-center">
          <Icon name="mdi-folder-open-outline" class="h-12 w-12 mx-auto mb-2" />
          <p>This directory is empty</p>
        </div>
      </div>

      <!-- File/directory grid -->
      <div v-else class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-3 overflow-y-auto flex-1 p-1">
        <!-- Directories (always all visible, navigation elements) -->
        <button
          v-for="entry in dirs"
          :key="entry.path"
          class="flex flex-col items-center gap-1 p-2 rounded-lg hover:bg-base-200 transition-colors text-center cursor-pointer border border-transparent hover:border-base-300"
          @click="onEntryClick(entry)"
        >
          <Icon name="mdi-folder" class="h-10 w-10 text-yellow-500" />
          <span class="text-xs break-all line-clamp-2 leading-tight">{{ entry.name }}</span>
        </button>

        <!-- Files (paginated, infinite scroll) -->
        <button
          v-for="entry in files"
          :key="entry.path"
          class="flex flex-col items-center gap-1 p-2 rounded-lg hover:bg-base-200 transition-colors text-center cursor-pointer border border-transparent hover:border-base-300"
          :class="{ 'opacity-60': !entry.isImage }"
          @click="onEntryClick(entry)"
        >
          <!-- Image thumbnail -->
          <img
            v-if="entry.isImage"
            :src="thumbSrc(entry.path)"
            class="h-16 w-full object-cover rounded bg-base-300"
            loading="lazy"
            decoding="async"
            :alt="entry.name"
          />

          <!-- Generic file icon -->
          <Icon v-else name="mdi-file-outline" class="h-10 w-10 text-gray-400" />

          <span class="text-xs break-all line-clamp-2 leading-tight">{{ entry.name }}</span>
          <span class="text-xs text-gray-400">{{ formatSize(entry.size) }}</span>
        </button>

        <!-- Sentinel element for infinite scroll trigger -->
        <div v-if="hasMore" ref="sentinel" class="col-span-full flex justify-center py-4">
          <span v-if="loadingMore" class="loading loading-spinner loading-sm"></span>
        </div>
      </div>
    </div>

    <template v-if="!loading && !errorMsg" #footer>
      <div class="text-xs text-gray-400">
        {{ dirs.length + files.length }} item{{ (dirs.length + files.length) !== 1 ? 's' : '' }} ({{ totalFiles }} file{{ totalFiles !== 1 ? 's' : '' }} total)
      </div>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import type { FileEntry } from "~~/lib/api/classes/import";

type Props = {
  modelValue: boolean;
};

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
});

const emit = defineEmits<{
  (e: "update:modelValue", v: boolean): void;
  (e: "select", entry: FileEntry): void;
}>();

const open = useVModel(props, "modelValue", emit);

const api = useUserApi();
const toast = useNotifier();

const loading = ref(false);
const errorMsg = ref("");
const dirs = ref<FileEntry[]>([]);
const files = ref<FileEntry[]>([]);
const currentPath = ref("");
const currentPage = ref(1);
const totalFiles = ref(0);
const pageSize = 50;
const loadingMore = ref(false);
const sentinel = ref<HTMLElement | null>(null);

const hasMore = computed(() => files.value.length < totalFiles.value);

// Breadcrumbs: split path segments
const breadcrumbs = computed(() => {
  if (!currentPath.value) return [];
  return currentPath.value.split("/").filter(Boolean);
});

function thumbSrc(filePath: string): string {
  return api.items.importDir.thumbURL(filePath, 120);
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

async function loadDir(subPath: string) {
  loading.value = true;
  errorMsg.value = "";
  currentPage.value = 1;
  dirs.value = [];
  files.value = [];
  totalFiles.value = 0;

  const { data, error } = await api.items.importDir.browse(subPath, 1, pageSize);

  loading.value = false;

  if (error) {
    errorMsg.value = `Failed to browse directory: ${error}`;
    return;
  }

  if (data) {
    dirs.value = data.dirs || [];
    files.value = data.files || [];
    totalFiles.value = data.total || 0;
  }
}

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return;
  loadingMore.value = true;
  const nextPage = currentPage.value + 1;

  const { data, error } = await api.items.importDir.browse(currentPath.value, nextPage, pageSize);

  loadingMore.value = false;

  if (error) return;

  if (data && data.files) {
    files.value.push(...data.files);
    currentPage.value = data.page;
    totalFiles.value = data.total;
  }
}

// Infinite scroll via IntersectionObserver — auto-cleaned by vueuse when sentinel unmounts
useIntersectionObserver(
  sentinel,
  ([{ isIntersecting }]) => {
    if (isIntersecting && hasMore.value && !loading.value && !loadingMore.value) {
      loadMore();
    }
  },
  { rootMargin: '200px' }
);

function goToDir(index: number) {
  if (index < 0) {
    currentPath.value = "";
  } else {
    const crumbs = breadcrumbs.value.slice(0, index + 1);
    currentPath.value = crumbs.join("/");
  }
  loadDir(currentPath.value);
}

function onEntryClick(entry: FileEntry) {
  if (entry.isDir) {
    // Navigate into directory — append to current path
    currentPath.value = currentPath.value
      ? `${currentPath.value}/${entry.name}`
      : entry.name;
    loadDir(currentPath.value);
  } else {
    // Import the file
    emit("select", entry);
  }
}

// Reload on open
watch(open, async (isOpen) => {
  if (isOpen) {
    currentPath.value = "";
    await loadDir("");
  }
});
</script>
