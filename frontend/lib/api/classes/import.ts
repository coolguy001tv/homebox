import { BaseAPI, route } from "../base";
import type { QueryValue } from "../base";
import { AttachmentTypes } from "../types/non-generated";
import { ItemOut } from "../types/data-contracts";

export interface FileEntry {
  name: string;
  path: string;
  size: number;
  isDir: boolean;
  extension: string;
  modified: string;
  isImage: boolean;
}

export interface BrowseResult {
  dirs: FileEntry[];
  files: FileEntry[];
  page: number;
  pageSize: number;
  total: number;
}

export class ImportAPI extends BaseAPI {
  browse(path: string = "", page: number = 1, pageSize: number = 50) {
    const params: Record<string, QueryValue> = { page, pageSize };
    if (path) params.path = path;
    return this.http.get<BrowseResult>({ url: route("/import/browse", params) });
  }

  importToItem(id: string, sourcePath: string, type: AttachmentTypes | null = null) {
    return this.http.post<{ sourcePath: string; type?: string }, ItemOut>({
      url: route(`/items/${id}/attachments/import`),
      body: { sourcePath, type: type || undefined },
    });
  }

  thumbURL(importPath: string, width: number): string {
    const encoded = encodeURIComponent(importPath);
    // Build the URL manually to avoid route() double-prefixing.
    // route() prepends /api/v1, and authURL() calls route() again,
    // so we construct the full path here directly.
    const path = `/import/thumb?path=${encoded}&w=${width}`;
    return this.authURL(path);
  }
}
