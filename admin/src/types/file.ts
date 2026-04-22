// 文件信息
export interface FileInfo {
  id: number;
  filename: string;
  original_name: string;
  file_url: string;
  file_type: string;
  file_size: number;
  upload_type: string;
  upload_time: string;
  status: number; // 0:未使用 1:使用中
}

// 文件列表查询
export interface FileListQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  file_type?: string;
  status?: number;
  upload_type?: string;
  min_size?: number;
  max_size?: number;
  start_time?: string;
  end_time?: string;
}

// 文件列表响应
export interface FileListData {
  list: FileInfo[];
  total: number;
  page: number;
  page_size: number;
}
