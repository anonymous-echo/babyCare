/**
 * 用户信息接口 (去家庭化架构)
 */
export interface UserInfo {
  openid: string
  nickName: string
  avatarUrl: string
  defaultBabyId?: string // 默认宝宝ID
  babies?: BabyProfile[] // 用户可访问的宝宝列表
  createTime: number
}

/**
 * 宝宝档案接口 (去家庭化架构)
 */
export interface BabyProfile {
  babyId: string  // 改为 babyId (与后端一致)
  name: string
  nickname?: string
  birthDate: string // YYYY-MM-DD 格式
  gender: 'male' | 'female'
  avatarUrl?: string
  creatorId: string // 创建者 openid
  familyGroup?: string // 可选的家庭分组名称
  createTime: number
  updateTime: number
}

/**
 * 宝宝协作者角色
 */
export type CollaboratorRole = 'admin' | 'editor' | 'viewer'

/**
 * 访问类型
 */
export type AccessType = 'permanent' | 'temporary'

/**
 * 宝宝协作者接口 (替代 FamilyMember)
 */
export interface BabyCollaborator {
  openid: string
  nickName: string
  avatarUrl: string
  role: CollaboratorRole
  relationship?: string // 与宝宝的关系: 爸爸, 妈妈, 爷爷, 奶奶, 外公, 外婆等
  accessType: AccessType
  expiresAt?: number // 临时权限过期时间
  joinTime: number
}

/**
 * 喂养类型
 */
export type FeedingType = 'breast' | 'bottle' | 'food'

/**
 * 喂养侧
 */
export type BreastSide = 'left' | 'right' | 'both'

/**
 * 单次喂养会话(多段式记录)
 */
export interface FeedingSession {
  side: BreastSide
  startTime: number   // 开始时间戳(毫秒)
  endTime?: number    // 结束时间戳(毫秒,可选)
  duration: number    // 时长(秒)
}

/**
 * 母乳喂养记录
 */
export interface BreastFeeding {
  type: 'breast'
  side: BreastSide // 主要喂养侧(向后兼容)
  duration: number // 总时长(秒)
  leftDuration?: number // 左侧总时长(秒)
  rightDuration?: number // 右侧总时长(秒)
  sessions?: FeedingSession[] // 多段式记录(可选,新功能)
}

/**
 * 奶瓶喂养记录
 */
export interface BottleFeeding {
  type: 'bottle'
  bottleType: 'formula' | 'breast-milk' // 配方奶或母乳
  amount: number // 奶量(ml)
  unit: 'ml' | 'oz'
  remaining?: number // 剩余量
}

/**
 * 辅食记录
 */
export interface FoodFeeding {
  type: 'food'
  foodName: string // 辅食名称
  note?: string // 备注(接受程度、过敏反应等)
}

/**
 * 喂养记录联合类型
 */
export type FeedingDetail = BreastFeeding | BottleFeeding | FoodFeeding

/**
 * 喂养提醒设置
 */
export interface FeedingReminderSetting {
  enabled: boolean          // 是否启用提醒
  intervalMinutes: number   // 间隔(分钟)
  nextReminderTime: number  // 下次提醒时间戳(毫秒)
}

/**
 * 喂养提醒用户偏好
 */
export interface FeedingReminderPreferences {
  breast: number  // 母乳默认间隔(分钟)
  bottle: number  // 奶瓶默认间隔(分钟)
  food: number    // 辅食默认间隔(分钟)
}

/**
 * 喂养记录
 */
export interface FeedingRecord {
  id: string
  babyId: string
  time: number // 时间戳
  detail: FeedingDetail
  createBy: string // 创建人 openid
  createByName: string // 冗余:创建者昵称
  createByAvatar: string // 冗余:创建者头像
  createTime: number
  reminderSetting?: FeedingReminderSetting // 提醒设置(可选)
}

/**
 * 排泄类型
 */
export type DiaperType = 'pee' | 'poop' | 'both'

/**
 * 大便颜色
 */
export type PoopColor = 'yellow' | 'green' | 'brown' | 'black' | 'red' | 'white'

/**
 * 大便性状
 */
export type PoopTexture = 'watery' | 'loose' | 'paste' | 'soft' | 'formed' | 'hard'

/**
 * 排泄记录
 */
export interface DiaperRecord {
  id: string
  babyId: string
  time: number
  type: DiaperType
  poopColor?: PoopColor // 大便颜色(仅在有大便时)
  poopTexture?: PoopTexture // 大便性状(仅在有大便时)
  note?: string
  createBy: string
  createByName: string // 冗余:创建者昵称
  createByAvatar: string // 冗余:创建者头像
  createTime: number
}

/**
 * 睡眠类型
 */
export type SleepType = 'nap' | 'night'

/**
 * 睡眠记录
 */
export interface SleepRecord {
  id: string
  babyId: string
  startTime: number // 开始时间戳
  endTime?: number // 结束时间戳(进行中时为空)
  duration?: number // 时长(分钟)
  type: SleepType
  createBy: string
  createByName: string // 冗余:创建者昵称
  createByAvatar: string // 冗余:创建者头像
  createTime: number
}

/**
 * 其他记录类型
 */
export type OtherRecordType = 'spit' | 'pump' | 'growth' | 'medicine' | 'temperature' | 'bath' | 'milestone'

/**
 * 吐奶记录
 */
export interface SpitRecord {
  type: 'spit'
  time: number
}

/**
 * 泵奶记录
 */
export interface PumpRecord {
  type: 'pump'
  duration: number // 时长(分钟)
  amount: number // 产量(ml)
}

/**
 * 生长记录
 */
export interface GrowthRecord {
  id: string
  babyId: string
  time: number
  height?: number // 身高(cm)
  weight?: number // 体重(kg)
  headCircumference?: number // 头围(cm)
  createBy: string
  createByName: string // 冗余:创建者昵称
  createByAvatar: string // 冗余:创建者头像
  createTime: number
}

/**
 * 其他事件记录
 */
export interface OtherEventRecord {
  id: string
  babyId: string
  time: number
  recordType: OtherRecordType
  detail: SpitRecord | PumpRecord | GrowthRecord | {
    type: OtherRecordType
    note: string
  }
  createBy: string
  createByName: string // 冗余:创建者昵称
  createByAvatar: string // 冗余:创建者头像
  createTime: number
}

/**
 * 统一的记录类型
 */
export type RecordType = 'feeding' | 'diaper' | 'sleep' | 'other'

/**
 * 所有记录的联合类型
 */
export type Record = FeedingRecord | DiaperRecord | SleepRecord | OtherEventRecord

/**
 * 家庭成员角色
 * @deprecated 此类型已弃用,请使用 CollaboratorRole 替代
 */
export type FamilyRole = 'admin' | 'member'

/**
 * 家庭成员
 * @deprecated 此类型已弃用,请使用 BabyCollaborator 替代
 */
export interface FamilyMember {
  openid: string
  nickName: string
  avatarUrl: string
  role: FamilyRole
  joinTime: number
}

/**
 * 家庭信息
 * @deprecated 此类型已弃用,去家庭化架构已移除家庭概念
 */
export interface FamilyInfo {
  familyId: string
  familyName: string
  description?: string
  creatorId: string
  createTime: number
  members: FamilyMember[]
  babyIds: string[]  // 关联的宝宝ID列表
}

/**
 * 邀请码信息
 * @deprecated 此类型已弃用,请使用基于宝宝的邀请机制
 */
export interface InvitationInfo {
  invitationCode: string
  familyId: string
  familyName: string
  creatorName: string
  qrCodeUrl?: string
  shareUrl?: string
  expiresAt: number
  createTime: number
}

/**
 * 同步状态
 */
export type SyncStatus = 'idle' | 'syncing' | 'success' | 'error'

/**
 * 同步配置
 */
export interface SyncConfig {
  autoSync: boolean       // 是否自动同步
  syncInterval: number    // 同步间隔(毫秒)
  wifiOnly: boolean       // 仅Wi-Fi同步
  lastSyncTime: number    // 上次同步时间
}

/**
 * 疫苗类型
 */
export type VaccineType =
  | 'BCG'              // 卡介苗
  | 'HepB'             // 乙肝疫苗
  | 'OPV'              // 脊灰疫苗(口服)
  | 'IPV'              // 脊灰疫苗(注射)
  | 'DTaP'             // 百白破疫苗
  | 'MR'               // 麻风疫苗
  | 'MMR'              // 麻腮风疫苗
  | 'JE'               // 乙脑疫苗
  | 'MeningB'          // 流脑B疫苗
  | 'MeningAC'         // 流脑AC疫苗
  | 'HepA'             // 甲肝疫苗
  | 'Varicella'        // 水痘疫苗
  | 'Pneumococcal'     // 肺炎疫苗
  | 'Rotavirus'        // 轮状病毒疫苗
  | 'Hib'              // Hib疫苗
  | 'Influenza'        // 流感疫苗
  | 'Other'            // 其他

/**
 * 疫苗计划项
 */
export interface VaccinePlan {
  id: string
  vaccineType: VaccineType
  vaccineName: string      // 疫苗名称
  description?: string     // 说明
  ageInMonths: number      // 接种月龄
  doseNumber: number       // 第几针(1, 2, 3...)
  isRequired: boolean      // 是否必打
  reminderDays: number     // 提前提醒天数
}

/**
 * 疫苗接种记录
 */
export interface VaccineRecord {
  id: string
  babyId: string
  planId: string           // 关联的计划ID
  vaccineType: VaccineType
  vaccineName: string
  doseNumber: number
  vaccineDate: number      // 接种日期(时间戳)
  hospital?: string        // 接种医院
  batchNumber?: string     // 疫苗批号
  doctor?: string          // 接种医生
  reaction?: string        // 不良反应
  note?: string           // 备注
  createBy: string
  createByName: string // 冗余:创建者昵称
  createByAvatar: string // 冗余:创建者头像
  createTime: number
}

/**
 * 疫苗提醒状态
 */
export type VaccineReminderStatus = 'upcoming' | 'due' | 'overdue' | 'completed'

/**
 * 疫苗提醒
 */
export interface VaccineReminder {
  id: string
  babyId: string
  planId: string
  vaccineName: string
  doseNumber: number
  scheduledDate: number    // 预定接种日期
  status: VaccineReminderStatus
  reminderSent: boolean    // 是否已发送提醒
  createTime: number
}

/**
 * 订阅消息模板类型
 */
export type SubscribeMessageType =
  | 'vaccine_reminder'       // 疫苗接种提醒
  | 'breast_feeding_reminder' // 母乳喂养提醒
  | 'bottle_feeding_reminder' // 奶瓶喂养提醒
  | 'pump_reminder'          // 吸奶器吸奶提醒
  | 'feeding_duration_alert' // 喂奶时间过长提醒

/**
 * 订阅消息模板配置
 */
export interface SubscribeMessageTemplate {
  type: SubscribeMessageType
  templateId: string         // 微信模板ID
  title: string              // 模板标题
  keywords: string[]         // 关键词列表
  description?: string       // 模板说明
  icon?: string              // 模板图标(本地路径)
  priority: number           // 优先级(1-5,5最高)
}

/**
 * 订阅消息授权状态
 */
export type SubscribeAuthStatus =
  | 'unknown'     // 未知(未申请过)
  | 'accept'      // 用户同意
  | 'reject'      // 用户拒绝
  | 'ban'         // 用户永久拒绝(拒绝3次后)

/**
 * 订阅消息授权记录
 */
export interface SubscribeAuthRecord {
  type: SubscribeMessageType
  templateId: string
  status: SubscribeAuthStatus
  lastRequestTime: number    // 上次申请时间
  requestCount: number       // 累计申请次数
  rejectCount: number        // 累计拒绝次数
  acceptCount: number        // 累计同意次数
  updateTime: number
}

/**
 * 订阅消息提醒配置
 */
export interface SubscribeReminderConfig {
  type: SubscribeMessageType
  enabled: boolean           // 是否启用提醒
  advanceDays?: number       // 提前天数(疫苗提醒)
  intervalMinutes?: number   // 间隔分钟数(喂养提醒)
  startTime?: string         // 提醒开始时间 HH:mm
  endTime?: string           // 提醒结束时间 HH:mm
}

/**
 * 订阅消息引导显示记录
 */
export interface SubscribeGuideRecord {
  type: SubscribeMessageType
  showCount: number          // 显示次数
  lastShowTime: number       // 上次显示时间
  nextShowTime: number       // 下次可显示时间
  isDismissedForever: boolean // 是否永久不再显示
}

/**
 * API 响应基础接口
 */
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  timestamp?: number
}