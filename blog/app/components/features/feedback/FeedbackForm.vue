<script setup lang="ts">
import { submitFeedback, getFeedbackByTicketNo } from '@/composables/api/feedback';
import type { SubmitFeedbackParams, ReportType, Feedback } from '@@/types/feedback';

interface Emits {
  (e: 'success'): void;
}

const emit = defineEmits<Emits>();
const { success, error: errorToast, warning } = useToast();

// 提交成功状态
const submitSuccess = ref(false);
const submittedTicketNo = ref('');

// 工单查询对话框
const showQueryDialog = ref(false);
const queryTicketNo = ref('');
const queryLoading = ref(false);
const queryResult = ref<Feedback | null>(null);

// 表单数据
const formData = reactive({
  reportUrl: '',
  reportType: '' as ReportType | '',
  email: '',
  // 版权投诉 - 侵权说明
  copyrightDescription: '',
  // 不当内容投诉 - 投诉内容
  inappropriateDescription: '',
  // 摘要问题反馈 - 反馈内容
  summaryDescription: '',
  // 功能建议反馈 - 功能描述
  suggestionDescription: '',
  // 功能建议反馈 - 使用场景
  suggestionReason: '',
  // 版权投诉 - 权利人证明
  copyrightProofFiles: [] as File[],
  // 版权投诉 - 侵权内容证明
  copyrightInfringementFiles: [] as File[],
  // 不当内容 - 证据截图
  inappropriateEvidenceFiles: [] as File[],
  // 摘要问题 - 相关截图
  summaryScreenshotFiles: [] as File[],
  // 功能建议 - 相关附件
  suggestionAttachmentFiles: [] as File[],
  // 不适当内容投诉的多选原因
  inappropriateReasons: [] as string[],
  // 摘要问题类型
  summaryIssueType: '',
});

// 提交状态
const submitting = ref(false);

// 反馈类型选项
const feedbackTypes = [
  { value: 'copyright', label: '版权侵权内容投诉' },
  { value: 'inappropriate', label: '不当内容举报投诉' },
  { value: 'summary', label: '文章摘要问题反馈' },
  { value: 'suggestion', label: '功能建议优化反馈' },
];

// 不适合内容的投诉原因选项
const inappropriateReasons = [
  '欺诈',
  '色情低俗',
  '诽谤',
  '传播不实信息',
  '违法犯罪',
  '骚扰',
  '诱导',
  '混淆他人',
  '恶意营销',
  '隐私侵权收集',
];

// 文章摘要问题类型
const summaryIssueTypes = [
  { value: 'inappropriate', label: '生成的内容包含恶意内容' },
  { value: 'mismatch', label: '生成内容与文章不符' },
];

// 文件上传引用 - 每个上传器独立的ref
const copyrightProofInputRef = ref<HTMLInputElement>();
const copyrightInfringementInputRef = ref<HTMLInputElement>();
const inappropriateEvidenceInputRef = ref<HTMLInputElement>();
const summaryScreenshotInputRef = ref<HTMLInputElement>();
const suggestionAttachmentInputRef = ref<HTMLInputElement>();

// 文件数组映射
type FileArrayKey =
  | 'copyrightProofFiles'
  | 'copyrightInfringementFiles'
  | 'inappropriateEvidenceFiles'
  | 'summaryScreenshotFiles'
  | 'suggestionAttachmentFiles';

const fileArrayMap: Record<string, FileArrayKey> = {
  copyrightProof: 'copyrightProofFiles',
  copyrightInfringement: 'copyrightInfringementFiles',
  inappropriateEvidence: 'inappropriateEvidenceFiles',
  summaryScreenshot: 'summaryScreenshotFiles',
  suggestionAttachment: 'suggestionAttachmentFiles',
};

// 处理文件选择
const handleFileUpload = (event: Event, targetKey: keyof typeof fileArrayMap) => {
  const target = event.target as HTMLInputElement;
  const files = Array.from(target.files || []);

  if (files.length === 0) return;

  // 验证每个文件
  const validFiles: File[] = [];
  for (const file of files) {
    const validationError = validateFile(file, '反馈投诉');
    if (validationError) {
      warning(validationError);
      continue;
    }
    validFiles.push(file);
  }

  if (validFiles.length > 0) {
    const arrayKey = fileArrayMap[targetKey] as FileArrayKey;
    formData[arrayKey] = [...formData[arrayKey], ...validFiles];
  }

  // 清空input以允许重复选择同一文件
  target.value = '';
};

// 移除文件
const removeFile = (index: number, targetKey: keyof typeof fileArrayMap) => {
  const arrayKey = fileArrayMap[targetKey] as FileArrayKey;
  formData[arrayKey].splice(index, 1);
};

// 验证表单
const validateForm = (): boolean => {
  // 检查公共必填项
  if (!formData.reportUrl.trim() || !formData.reportType) {
    warning('请完整填写表单必填项');
    return false;
  }

  // 检查各类型独立的必填项
  if (formData.reportType === 'copyright') {
    if (!formData.copyrightDescription.trim()) {
      warning('请完整填写表单必填项');
      return false;
    }
  } else if (formData.reportType === 'inappropriate') {
    if (formData.inappropriateReasons.length === 0 || !formData.inappropriateDescription.trim()) {
      warning('请完整填写表单必填项');
      return false;
    }
  } else if (formData.reportType === 'summary') {
    if (!formData.summaryIssueType || !formData.summaryDescription.trim()) {
      warning('请完整填写表单必填项');
      return false;
    }
  } else if (formData.reportType === 'suggestion') {
    if (!formData.suggestionDescription.trim() || !formData.suggestionReason.trim()) {
      warning('请完整填写表单必填项');
      return false;
    }
  }

  return true;
};

// 上传所有文件并返回URL数组
const uploadFiles = async (files: File[]): Promise<string[]> => {
  if (files.length === 0) return [];

  const uploadPromises = files.map(file => uploadFile(file, '反馈投诉'));
  const uploadResults = await Promise.all(uploadPromises);
  return uploadResults.map(result => result.file_url);
};

// 提交表单
const handleSubmit = async () => {
  if (!validateForm()) return;

  submitting.value = true;
  try {
    // 收集当前类型的文件并上传
    let allFiles: File[] = [];
    let description = '';
    let reason = '';

    switch (formData.reportType) {
      case 'copyright':
        allFiles = [...formData.copyrightProofFiles, ...formData.copyrightInfringementFiles];
        description = formData.copyrightDescription;
        // 版权侵权需要记录文件分组信息，存在reason字段
        reason = JSON.stringify({
          proofCount: formData.copyrightProofFiles.length,
          infringementCount: formData.copyrightInfringementFiles.length,
        });
        break;
      case 'inappropriate':
        allFiles = [...formData.inappropriateEvidenceFiles];
        description = formData.inappropriateDescription;
        reason = formData.inappropriateReasons.join(', ');
        break;
      case 'summary':
        allFiles = [...formData.summaryScreenshotFiles];
        description = formData.summaryDescription;
        // 提交中文标签而不是英文key，保持系统统一
        reason =
          summaryIssueTypes.find(item => item.value === formData.summaryIssueType)?.label ||
          formData.summaryIssueType;
        break;
      case 'suggestion':
        allFiles = [...formData.suggestionAttachmentFiles];
        description = formData.suggestionDescription;
        reason = formData.suggestionReason;
        break;
    }

    const attachmentUrls = await uploadFiles(allFiles);

    const submitData: SubmitFeedbackParams = {
      reportUrl: formData.reportUrl,
      reportType: formData.reportType as ReportType,
      email: formData.email || undefined,
      description: description || '-', // 版权投诉没有描述字段，用占位符
      reason: reason || undefined,
      attachmentFiles: attachmentUrls.length > 0 ? attachmentUrls : undefined,
    };

    const result = await submitFeedback(submitData);
    submittedTicketNo.value = result.ticket_no;
    submitSuccess.value = true;
    resetForm();
    emit('success');
  } catch (error: unknown) {
    const err = error as Error;
    errorToast(err.message || '提交失败，请稍后重试');
  } finally {
    submitting.value = false;
  }
};

// 重置表单
const resetForm = () => {
  formData.reportUrl = '';
  formData.reportType = '' as ReportType | '';
  formData.email = '';
  formData.copyrightDescription = '';
  formData.inappropriateDescription = '';
  formData.summaryDescription = '';
  formData.suggestionDescription = '';
  formData.suggestionReason = '';
  formData.copyrightProofFiles = [];
  formData.copyrightInfringementFiles = [];
  formData.inappropriateEvidenceFiles = [];
  formData.summaryScreenshotFiles = [];
  formData.suggestionAttachmentFiles = [];
  formData.inappropriateReasons = [];
  formData.summaryIssueType = '';
};

// 处理不适合内容原因的选择
const toggleReason = (reason: string) => {
  const index = formData.inappropriateReasons.indexOf(reason);
  if (index > -1) {
    formData.inappropriateReasons.splice(index, 1);
  } else {
    formData.inappropriateReasons.push(reason);
  }
};

// 返回表单
const backToForm = () => {
  submitSuccess.value = false;
  submittedTicketNo.value = '';
};

// 打开工单查询对话框
const openQueryDialog = () => {
  showQueryDialog.value = true;
  queryTicketNo.value = '';
  queryResult.value = null;
};

// 关闭工单查询对话框
const closeQueryDialog = () => {
  showQueryDialog.value = false;
  queryTicketNo.value = '';
  queryResult.value = null;
};

// 查询工单
const handleQueryTicket = async () => {
  if (!queryTicketNo.value.trim()) {
    warning('请输入工单号');
    return;
  }

  queryLoading.value = true;
  try {
    queryResult.value = await getFeedbackByTicketNo(queryTicketNo.value.trim());
    success('查询成功');
  } catch (error: unknown) {
    const err = error as Error;
    errorToast(err.message || '未找到该工单，请检查工单号是否正确');
    queryResult.value = null;
  } finally {
    queryLoading.value = false;
  }
};

// 获取状态标签文本
const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '待处理',
    resolved: '已解决',
    closed: '已关闭',
  };
  return labels[status] || status;
};

// 格式化日期
const formatDate = (date?: string) => {
  if (!date) return '未回复';
  return new Date(date).toLocaleString('zh-CN');
};
</script>

<template>
  <div class="feedback-form">
    <!-- 成功状态 -->
    <div v-if="submitSuccess" class="success-state">
      <div class="success-icon">
        <i class="ri-checkbox-circle-fill" />
      </div>
      <h3 class="success-title">提交成功！我们会尽快处理</h3>
      <div class="ticket-no">{{ submittedTicketNo }}</div>
      <p class="success-notice">请妥善保存工单号，这将作为结果查询的唯一方式</p>
      <button type="button" class="back-btn" @click="backToForm">
        <i class="ri-arrow-left-line" />
        返回
      </button>
    </div>

    <!-- 表单 -->
    <form v-else @submit.prevent="handleSubmit">
      <!-- 01 反馈的地址 -->
      <div class="form-group">
        <label class="label required">
          <span class="number">*01</span>
          请输入反馈的地址
        </label>
        <input
          v-model="formData.reportUrl"
          type="url"
          class="input"
          placeholder="请输入需要反馈的网页地址"
          :disabled="submitting"
        />
      </div>

      <!-- 02 反馈类型 -->
      <div class="form-group">
        <label class="label required">
          <span class="number">*02</span>
          反馈类型
        </label>
        <div class="radio-group">
          <label
            v-for="option in feedbackTypes"
            :key="option.value"
            class="radio-item"
            :class="{ active: formData.reportType === option.value }"
          >
            <input
              v-model="formData.reportType"
              type="radio"
              :value="option.value"
              :disabled="submitting"
            />
            <span class="radio-label">{{ option.label }}</span>
          </label>
        </div>
      </div>

      <!-- 03 联系邮箱（可选） -->
      <div class="form-group">
        <label class="label">
          <span class="number">03</span>
          联系邮箱
        </label>
        <div class="info-text">提供邮箱可以方便我们与您联系处理结果</div>
        <input
          v-model="formData.email"
          type="email"
          class="input"
          placeholder="请输入邮箱地址"
          :disabled="submitting"
        />
      </div>

      <!-- 内容版权侵权投诉 -->
      <template v-if="formData.reportType === 'copyright'">
        <div class="form-group">
          <label class="label required">
            <span class="number">*04</span>
            请输入侵权说明
          </label>
          <textarea
            v-model="formData.copyrightDescription"
            class="textarea"
            placeholder="请详细说明侵权情况"
            rows="6"
            :disabled="submitting"
          />
        </div>

        <div class="form-group">
          <label class="label required">
            <span class="number">*05</span>
            声明您是被侵权的权利人
          </label>
          <div class="info-text">
            可通过提供不限于以下的形式：商标注册书、专利著作权证明、后台截图等
          </div>
          <div class="file-upload-section">
            <input
              ref="copyrightProofInputRef"
              type="file"
              multiple
              accept="image/*,.pdf,.doc,.docx"
              style="display: none"
              @change="handleFileUpload($event, 'copyrightProof')"
            />
            <div class="upload-area" @click="copyrightProofInputRef?.click()">
              <i class="ri-add-line" />
              <span>点击选择<br />图片/文件</span>
            </div>
          </div>
          <!-- 文件列表 -->
          <div v-if="formData.copyrightProofFiles.length > 0" class="uploaded-files-inline">
            <div class="file-list">
              <div
                v-for="(file, index) in formData.copyrightProofFiles"
                :key="index"
                class="file-item"
              >
                <i class="ri-file-line" />
                <span class="file-name">{{ file.name }}</span>
                <button
                  type="button"
                  class="remove-btn"
                  @click="removeFile(index, 'copyrightProof')"
                >
                  <i class="ri-close-line" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="form-group">
          <label class="label required">
            <span class="number">*06</span>
            声明侵权的内容
          </label>
          <div class="info-text">请提供具有时间依据的内容对比材料</div>
          <div class="file-upload-section">
            <input
              ref="copyrightInfringementInputRef"
              type="file"
              multiple
              accept="image/*,.pdf,.doc,.docx"
              style="display: none"
              @change="handleFileUpload($event, 'copyrightInfringement')"
            />
            <div class="upload-area" @click="copyrightInfringementInputRef?.click()">
              <i class="ri-add-line" />
              <span>点击选择<br />图片/文件</span>
            </div>
          </div>
          <!-- 文件列表 -->
          <div v-if="formData.copyrightInfringementFiles.length > 0" class="uploaded-files-inline">
            <div class="file-list">
              <div
                v-for="(file, index) in formData.copyrightInfringementFiles"
                :key="index"
                class="file-item"
              >
                <i class="ri-file-line" />
                <span class="file-name">{{ file.name }}</span>
                <button
                  type="button"
                  class="remove-btn"
                  @click="removeFile(index, 'copyrightInfringement')"
                >
                  <i class="ri-close-line" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- 不适合的内容投诉 -->
      <template v-if="formData.reportType === 'inappropriate'">
        <div class="form-group">
          <label class="label required">
            <span class="number">*04</span>
            请选择投诉原因
          </label>
          <div class="checkbox-group">
            <label v-for="reason in inappropriateReasons" :key="reason" class="checkbox-item">
              <input
                type="checkbox"
                :checked="formData.inappropriateReasons.includes(reason)"
                :disabled="submitting"
                @change="toggleReason(reason)"
              />
              <span class="checkbox-label">{{ reason }}</span>
            </label>
          </div>
        </div>

        <div class="form-group">
          <label class="label required">
            <span class="number">*05</span>
            请输入投诉内容
          </label>
          <textarea
            v-model="formData.inappropriateDescription"
            class="textarea"
            placeholder="请输入"
            rows="6"
            :disabled="submitting"
          />
        </div>

        <div class="form-group">
          <label class="label">
            <span class="number">06</span>
            请上传证据截图
          </label>
          <div class="file-upload-section">
            <input
              ref="inappropriateEvidenceInputRef"
              type="file"
              multiple
              accept="image/*,.pdf,.doc,.docx"
              style="display: none"
              @change="handleFileUpload($event, 'inappropriateEvidence')"
            />
            <div class="upload-area" @click="inappropriateEvidenceInputRef?.click()">
              <i class="ri-add-line" />
              <span>点击选择<br />图片/文件</span>
            </div>
          </div>
          <!-- 文件列表 -->
          <div v-if="formData.inappropriateEvidenceFiles.length > 0" class="uploaded-files-inline">
            <div class="file-list">
              <div
                v-for="(file, index) in formData.inappropriateEvidenceFiles"
                :key="index"
                class="file-item"
              >
                <i class="ri-file-line" />
                <span class="file-name">{{ file.name }}</span>
                <button
                  type="button"
                  class="remove-btn"
                  @click="removeFile(index, 'inappropriateEvidence')"
                >
                  <i class="ri-close-line" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- 文章摘要问题反馈 -->
      <template v-if="formData.reportType === 'summary'">
        <div class="form-group">
          <label class="label required">
            <span class="number">*04</span>
            文章摘要反馈
          </label>
          <div class="info-text">
            文章摘要由AI模型自动生成，内容已经过人工审核。如发现生成内容有问题，请选择相应的反馈原因。
          </div>
          <div class="radio-group">
            <label
              v-for="issueType in summaryIssueTypes"
              :key="issueType.value"
              class="radio-item"
              :class="{ active: formData.summaryIssueType === issueType.value }"
            >
              <input
                v-model="formData.summaryIssueType"
                type="radio"
                :value="issueType.value"
                :disabled="submitting"
              />
              <span class="radio-label">{{ issueType.label }}</span>
            </label>
          </div>
        </div>

        <div class="form-group">
          <label class="label required">
            <span class="number">*05</span>
            请输入反馈内容
          </label>
          <textarea
            v-model="formData.summaryDescription"
            class="textarea"
            placeholder="请输入"
            rows="6"
            :disabled="submitting"
          />
        </div>

        <div class="form-group">
          <label class="label">
            <span class="number">06</span>
            请上传相关截图
          </label>
          <div class="file-upload-section">
            <input
              ref="summaryScreenshotInputRef"
              type="file"
              multiple
              accept="image/*,.pdf,.doc,.docx"
              style="display: none"
              @change="handleFileUpload($event, 'summaryScreenshot')"
            />
            <div class="upload-area" @click="summaryScreenshotInputRef?.click()">
              <i class="ri-add-line" />
              <span>点击选择<br />图片/文件</span>
            </div>
          </div>
          <!-- 文件列表 -->
          <div v-if="formData.summaryScreenshotFiles.length > 0" class="uploaded-files-inline">
            <div class="file-list">
              <div
                v-for="(file, index) in formData.summaryScreenshotFiles"
                :key="index"
                class="file-item"
              >
                <i class="ri-file-line" />
                <span class="file-name">{{ file.name }}</span>
                <button
                  type="button"
                  class="remove-btn"
                  @click="removeFile(index, 'summaryScreenshot')"
                >
                  <i class="ri-close-line" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- 功能建议反馈 -->
      <template v-if="formData.reportType === 'suggestion'">
        <div class="form-group">
          <label class="label required">
            <span class="number">*04</span>
            请输入功能描述
          </label>
          <textarea
            v-model="formData.suggestionDescription"
            class="textarea"
            placeholder="请详细描述您希望增加的功能"
            rows="4"
            :disabled="submitting"
          />
        </div>

        <div class="form-group">
          <label class="label required">
            <span class="number">*05</span>
            请输入使用场景
          </label>
          <textarea
            v-model="formData.suggestionReason"
            class="textarea"
            placeholder="请描述该功能的具体使用场景和解决的问题"
            rows="4"
            :disabled="submitting"
          />
        </div>

        <div class="form-group">
          <label class="label">
            <span class="number">06</span>
            请上传相关附件
          </label>
          <div class="info-text">可上传设计图、示意图等相关文件，帮助我们更好地理解您的建议</div>
          <div class="file-upload-section">
            <input
              ref="suggestionAttachmentInputRef"
              type="file"
              multiple
              accept="image/*,.pdf,.doc,.docx"
              style="display: none"
              @change="handleFileUpload($event, 'suggestionAttachment')"
            />
            <div class="upload-area" @click="suggestionAttachmentInputRef?.click()">
              <i class="ri-add-line" />
              <span>点击选择<br />图片/文件</span>
            </div>
          </div>
          <!-- 文件列表 -->
          <div v-if="formData.suggestionAttachmentFiles.length > 0" class="uploaded-files-inline">
            <div class="file-list">
              <div
                v-for="(file, index) in formData.suggestionAttachmentFiles"
                :key="index"
                class="file-item"
              >
                <i class="ri-file-line" />
                <span class="file-name">{{ file.name }}</span>
                <button
                  type="button"
                  class="remove-btn"
                  @click="removeFile(index, 'suggestionAttachment')"
                >
                  <i class="ri-close-line" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- 提交按钮 -->
      <div class="form-actions">
        <button type="submit" class="submit-btn" :disabled="submitting">
          <i class="ri-send-plane-fill" :class="{ 'ri-loader-4-line ri-spin': submitting }" />
          {{ submitting ? '提交中...' : '提交反馈' }}
        </button>
        <button type="button" class="query-btn" @click="openQueryDialog">
          <i class="ri-search-line" />
          工单查询
        </button>
      </div>
    </form>

    <!-- 工单查询对话框 -->
    <div v-if="showQueryDialog" class="dialog-overlay" @click.self="closeQueryDialog">
      <div class="dialog-content">
        <div class="dialog-header">
          <h3>工单查询</h3>
          <button class="close-btn" @click="closeQueryDialog">
            <i class="ri-close-line" />
          </button>
        </div>

        <div class="dialog-body">
          <div class="query-input-group">
            <input
              v-model="queryTicketNo"
              type="text"
              class="query-input"
              placeholder="请输入工单号，例如：FB20241108001"
              @keyup.enter="handleQueryTicket"
            />
            <button
              type="button"
              class="query-search-btn"
              :disabled="queryLoading"
              @click="handleQueryTicket"
            >
              <i class="ri-search-line" :class="{ 'ri-loader-4-line ri-spin': queryLoading }" />
              {{ queryLoading ? '查询中...' : '查询' }}
            </button>
          </div>

          <!-- 查询结果 -->
          <div v-if="queryResult" class="query-result">
            <div class="result-header">
              <span class="result-ticket">工单号：{{ queryResult.ticket_no }}</span>
              <span class="result-status" :class="`status-${queryResult.status}`">
                {{ getStatusLabel(queryResult.status) }}
              </span>
            </div>

            <div class="result-item">
              <label>反馈时间：</label>
              <span>{{ formatDate(queryResult.feedback_time) }}</span>
            </div>

            <div v-if="queryResult.admin_reply" class="result-item">
              <label>管理员回复：</label>
              <p class="reply-content">{{ queryResult.admin_reply }}</p>
            </div>

            <div v-if="queryResult.reply_time" class="result-item">
              <label>回复时间：</label>
              <span>{{ formatDate(queryResult.reply_time) }}</span>
            </div>

            <div v-if="!queryResult.admin_reply" class="result-empty">
              <i class="ri-time-line" />
              <p>您的反馈正在处理中，请耐心等待...</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

.feedback-form {
  max-width: 600px;
  margin: 0 auto;

  .form-group {
    margin-bottom: 1.5rem;

    .label {
      margin-bottom: 0.5rem;
      font-weight: 500;
      color: var(--font-color);

      .number {
        font-weight: 700;
      }
    }

    .info-text {
      margin-bottom: 0.5rem;
      font-size: 0.8rem;
      line-height: 1.4;
      color: var(--font-color);
      opacity: 0.6;
    }

    .input,
    .textarea {
      width: 100%;
      padding: 0.6rem;
      border: 1px solid var(--flec-border);
      border-radius: 4px;
      color: var(--font-color);
      background-color: transparent;

      &:focus {
        outline: none;
        border-color: var(--theme-color);
      }

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      &::placeholder {
        color: var(--font-color);
        opacity: 0.6;
      }
    }

    .textarea {
      resize: vertical;
      min-height: 80px;
      font-family: inherit;
      line-height: 1.5;
    }

    .radio-group,
    .checkbox-group {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .radio-item,
    .checkbox-item {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.6rem;
      border: 1px solid var(--flec-border);
      border-radius: 4px;
      cursor: pointer;

      &:hover {
        border-color: var(--theme-color);
      }

      &.active {
        border-color: var(--theme-color);
      }

      input[type='radio'],
      input[type='checkbox'] {
        margin: 0;
      }

      .radio-label,
      .checkbox-label {
        flex: 1;
        color: var(--font-color);
      }
    }

    .file-upload-section {
      .upload-area {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 1.5rem;
        border: 1px solid var(--flec-border);
        border-radius: 4px;
        background-color: transparent;
        cursor: pointer;
        text-align: center;

        &:hover:not(.uploading) {
          border-color: var(--flec-btn);
        }

        &.uploading {
          border-color: var(--theme-color);
          cursor: not-allowed;
          opacity: 0.8;
        }

        i {
          font-size: 1.5rem;
          color: var(--font-color);

          &.ri-spin {
            animation: spin 1s linear infinite;
            color: var(--theme-color);
          }
        }

        span {
          color: var(--font-color);
          font-size: 0.9rem;
          line-height: 1.3;
        }
      }
    }
  }

  .uploaded-files-inline {
    margin-top: 0.75rem;

    .file-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .file-item {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem;
      background-color: transparent;
      border: 1px solid var(--flec-border);
      border-radius: 4px;

      i {
        color: var(--flec-btn);
        font-size: 0.9rem;
      }

      .file-name {
        flex: 1;
        font-size: 0.9rem;
        color: var(--font-color);
      }

      .remove-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 20px;
        height: 20px;
        border: none;
        background: transparent;
        border-radius: 50%;
        color: var(--font-color);
        opacity: 0.7;
        cursor: pointer;

        &:hover {
          background-color: var(--flec-heavy-bg);
          opacity: 1;
        }

        i {
          font-size: 0.8rem;
        }
      }
    }
  }

  .form-actions {
    margin-top: 2rem;
    display: flex;
    gap: 1rem;

    .submit-btn,
    .query-btn {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      padding: 0.75rem;
      font-size: 1rem;
      border: none;
      border-radius: 4px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;

      i {
        font-size: 1rem;
      }
    }

    .submit-btn {
      background: var(--flec-btn);
      color: white;

      &:hover:not(:disabled) {
        background: var(--flec-btn-hover);
      }

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }
    }

    .query-btn {
      background: transparent;
      color: var(--flec-btn);
      border: 1px solid var(--flec-btn);

      &:hover {
        background: var(--flec-btn);
        color: white;
      }
    }
  }

  // 成功状态
  .success-state {
    text-align: center;
    padding: 3rem 2rem;

    .success-icon {
      margin-bottom: 1rem;

      i {
        font-size: 3.5rem;
        color: #52c41a;
      }
    }

    .success-title {
      font-size: 1.5rem;
      font-weight: 600;
      color: var(--font-color);
      margin-bottom: 1.5rem;
    }

    .ticket-no {
      font-size: 1.5rem;
      font-weight: 700;
      color: var(--font-color);
      letter-spacing: 2px;
      margin-bottom: 1rem;
    }

    .success-notice {
      font-size: 0.95rem;
      font-weight: 600;
      color: var(--font-color);
      margin-bottom: 2rem;
      line-height: 1.6;
    }

    .back-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem 2rem;
      background: var(--flec-btn);
      color: white;
      border: none;
      border-radius: 4px;
      font-size: 1rem;
      font-weight: 500;
      cursor: pointer;
      transition: background 0.2s;

      &:hover {
        background: var(--flec-btn-hover);
      }

      i {
        font-size: 1rem;
      }
    }
  }

  // 对话框样式
  .dialog-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    padding: 1rem;
  }

  .dialog-content {
    @extend .cardHover;
    width: 100%;
    max-width: 500px;
    max-height: 80vh;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .dialog-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--flec-border);

    h3 {
      font-size: 1.2rem;
      font-weight: 600;
      color: var(--font-color);
      margin: 0;
    }

    .close-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 32px;
      height: 32px;
      border: none;
      background: transparent;
      color: var(--font-color);
      opacity: 0.6;
      border-radius: 4px;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        background: var(--flec-heavy-bg);
        opacity: 1;
      }

      i {
        font-size: 1.2rem;
      }
    }
  }

  .dialog-body {
    padding: 1.5rem;
    overflow-y: auto;
  }

  .query-input-group {
    display: flex;
    gap: 0.75rem;
    margin-bottom: 1.5rem;

    .query-input {
      flex: 1;
      padding: 0.75rem;
      border: 1px solid var(--flec-border);
      border-radius: 4px;
      background: transparent;
      color: var(--font-color);
      font-size: 0.95rem;

      &:focus {
        outline: none;
        border-color: var(--flec-btn);
      }

      &::placeholder {
        color: var(--font-color);
        opacity: 0.5;
      }
    }

    .query-search-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem 1.5rem;
      background: var(--flec-btn);
      color: white;
      border: none;
      border-radius: 4px;
      font-size: 0.95rem;
      font-weight: 500;
      cursor: pointer;
      white-space: nowrap;
      transition: all 0.2s;

      &:hover:not(:disabled) {
        background: var(--flec-btn-hover);
      }

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      i {
        font-size: 1rem;
      }
    }
  }

  .query-result {
    background: var(--flec-light-bg);
    border: 1px solid var(--flec-border);
    border-radius: 8px;
    padding: 1.5rem;

    .result-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 1.5rem;
      padding-bottom: 1rem;
      border-bottom: 1px solid var(--flec-border);

      .result-ticket {
        font-size: 1rem;
        font-weight: 600;
        color: var(--font-color);
      }

      .result-status {
        padding: 0.25rem 0.75rem;
        border-radius: 12px;
        font-size: 0.85rem;
        font-weight: 500;

        &.status-pending {
          background: #fff3e0;
          color: #e65100;
        }

        &.status-resolved {
          background: #e8f5e9;
          color: #2e7d32;
        }

        &.status-closed {
          background: #e3f2fd;
          color: #1565c0;
        }
      }
    }

    .result-item {
      margin-bottom: 1rem;

      &:last-child {
        margin-bottom: 0;
      }

      label {
        display: block;
        font-size: 0.85rem;
        font-weight: 600;
        color: var(--font-color);
        opacity: 0.7;
        margin-bottom: 0.25rem;
      }

      span {
        font-size: 0.95rem;
        color: var(--font-color);
      }

      .reply-content {
        margin: 0.5rem 0 0 0;
        font-size: 0.95rem;
        line-height: 1.6;
        color: var(--font-color);
        white-space: pre-wrap;
      }
    }

    .result-empty {
      text-align: center;
      padding: 2rem 1rem;
      color: var(--font-color);
      opacity: 0.6;

      i {
        font-size: 3rem;
        margin-bottom: 1rem;
      }

      p {
        margin: 0;
        font-size: 0.95rem;
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .feedback-form {
    .form-group {
      margin-bottom: 1.25rem;

      .radio-item,
      .checkbox-item {
        padding: 0.5rem;
      }
    }

    .form-actions {
      flex-direction: column;
      gap: 0.75rem;
    }

    .success-state {
      padding: 2rem 1rem;

      .success-icon i {
        font-size: 3rem;
      }

      .success-title {
        font-size: 1.2rem;
      }

      .ticket-no {
        font-size: 1.2rem;
        letter-spacing: 1px;
      }

      .success-notice {
        font-size: 0.85rem;
      }
    }

    .dialog-content {
      margin: 0 1rem;
    }

    .query-input-group {
      flex-direction: column;
    }
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}
</style>
