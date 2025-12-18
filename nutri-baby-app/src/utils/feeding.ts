/**
 * 婴儿喂奶间隔和次数建议
 * 基于医学指南和育儿最佳实践
 */

/**
 * 喂奶间隔建议（以月龄分组）
 *
 * 参考资源:
 * - 美国儿科学会（AAP）
 * - 中国国家卫生健康委员会（NHC）
 * - 北京市卫生健康委员会
 * - 上海市卫生健康委员会
 *
 * 核心原则:
 * 1. 新生儿至3月龄: 按需喂养，逐步建立规律
 * 2. 4-6月龄: 定时喂养，夜间逐步减少
 * 3. 6-9月龄: 逐步加入辅食，夜间基本不喂
 * 4. 9-12月龄: 3-4次喂奶，辅食为主
 * 5. 12月龄以上: 规律三餐，乳品补充
 */

export interface FeedingGuideline {
  // 月龄范围
  minMonths: number;
  maxMonths: number;

  // 喂奶建议
  intervalMinHours: number;      // 最小间隔（小时）
  intervalMaxHours: number;      // 最大间隔（小时）
  recommendedTimesPerDay: number; // 建议日次数

  // 喂养特点
  feedingType: 'demand' | 'scheduled' | 'mixed'; // 按需/定时/混合

  // 添加辅食建议
  complementaryFoodsIntroduced: boolean;

  // 夜间喂奶建议
  nightFeeding: 'continue' | 'reduce' | 'stop'; // 继续/逐步减少/停止

  // 注意事项
  notes: string[];
}

/**
 * 喂奶指南数据库
 * 按月龄从小到大排列
 */
export const FEEDING_GUIDELINES: FeedingGuideline[] = [
  {
    minMonths: 0,
    maxMonths: 1,
    intervalMinHours: 2,
    intervalMaxHours: 3,
    recommendedTimesPerDay: 8,
    feedingType: 'demand',
    complementaryFoodsIntroduced: false,
    nightFeeding: 'continue',
    notes: [
      '新生儿按需喂养，胃容量小、消化快',
      '一般每2-3小时一次，每天8-12次',
      '观察宝宝饥饿信号和满足信号',
      '1个月左右可逐步形成间隔和规律',
    ],
  },
  {
    minMonths: 1,
    maxMonths: 3,
    intervalMinHours: 2.5,
    intervalMaxHours: 3.5,
    recommendedTimesPerDay: 7,
    feedingType: 'mixed',
    complementaryFoodsIntroduced: false,
    nightFeeding: 'continue',
    notes: [
      '从按需向定时逐步过渡',
      '婴儿开始建立喂养规律',
      '夜间哺乳仍需继续',
      '每个宝宝个体差异大，需灵活调整',
    ],
  },
  {
    minMonths: 3,
    maxMonths: 4,
    intervalMinHours: 3,
    intervalMaxHours: 4,
    recommendedTimesPerDay: 6,
    feedingType: 'scheduled',
    complementaryFoodsIntroduced: false,
    nightFeeding: 'reduce',
    notes: [
      '逐渐定时喂养，每3-4小时一次',
      '帮助婴儿建立较长的夜间睡眠',
      '开始考虑减少夜间哺乳次数',
      '宝宝可能开始显示对固体食物的兴趣',
    ],
  },
  {
    minMonths: 4,
    maxMonths: 6,
    intervalMinHours: 3,
    intervalMaxHours: 4,
    recommendedTimesPerDay: 6,
    feedingType: 'scheduled',
    complementaryFoodsIntroduced: true,
    nightFeeding: 'reduce',
    notes: [
      '定时喂养，每3-4小时一次',
      '可开始引入辅食（从6个月开始更正式）',
      '逐步减少夜间喂奶',
      '建议在医生指导下添加辅食',
    ],
  },
  {
    minMonths: 6,
    maxMonths: 9,
    intervalMinHours: 3.5,
    intervalMaxHours: 4.5,
    recommendedTimesPerDay: 5,
    feedingType: 'scheduled',
    complementaryFoodsIntroduced: true,
    nightFeeding: 'stop',
    notes: [
      '每日喂奶5次左右',
      '逐步加入辅食，部分进食成为"独立一餐"',
      '夜间喂奶逐渐减少至停止',
      '建立良好的进食习惯',
    ],
  },
  {
    minMonths: 9,
    maxMonths: 12,
    intervalMinHours: 4,
    intervalMaxHours: 5,
    recommendedTimesPerDay: 3,
    feedingType: 'scheduled',
    complementaryFoodsIntroduced: true,
    nightFeeding: 'stop',
    notes: [
      '每天喂奶3-4次为宜',
      '辅食比重逐步上升，成为主食',
      '夜间基本不再喂奶',
      '培养和家人一起进餐的习惯',
    ],
  },
  {
    minMonths: 12,
    maxMonths: 36,
    intervalMinHours: 4,
    intervalMaxHours: 6,
    recommendedTimesPerDay: 3,
    feedingType: 'scheduled',
    complementaryFoodsIntroduced: true,
    nightFeeding: 'stop',
    notes: [
      '主食为三餐，乳品补充',
      '每日乳品400-600毫升',
      '三餐两点的规律饮食',
      '培养与家庭一致的饮食习惯',
    ],
  },
];

/**
 * 根据月龄获取喂奶指南
 * @param ageInMonths 宝宝月龄（月）
 * @returns 对应的喂奶指南
 */
export function getFeedingGuidelineByAge(ageInMonths: number): FeedingGuideline {
  for (const guideline of FEEDING_GUIDELINES) {
    if (ageInMonths >= guideline.minMonths && ageInMonths < guideline.maxMonths + 1) {
      return guideline;
    }
  }
  // 如果超过所有范围，返回最后一个（12月龄以上）
  const lastGuideline = FEEDING_GUIDELINES[FEEDING_GUIDELINES.length - 1];
  if (!lastGuideline) {
    throw new Error('FEEDING_GUIDELINES is empty');
  }
  return lastGuideline;
}

/**
 * 计算推荐的喂奶间隔（分钟）
 *
 * @param ageInMonths 宝宝月龄
 * @param feedingType 当前喂养类型（可选，用于个性化调整）
 * @returns 推荐间隔的范围 { min: 分钟, max: 分钟, recommended: 分钟 }
 */
export function getRecommendedFeedingInterval(
  ageInMonths: number,
  feedingType?: 'breast' | 'bottle' | 'mixed'
) {
  const guideline = getFeedingGuidelineByAge(ageInMonths);

  const minMinutes = Math.floor(guideline.intervalMinHours * 60);
  const maxMinutes = Math.ceil(guideline.intervalMaxHours * 60);

  // 根据喂养类型微调
  let recommendedMinutes = Math.floor((minMinutes + maxMinutes) / 2);

  if (feedingType === 'breast') {
    // 母乳消化较快，可能需要更频繁的喂奶
    recommendedMinutes = Math.floor(recommendedMinutes * 0.9);
  } else if (feedingType === 'bottle') {
    // 配方奶消化较慢，可能能延长间隔
    recommendedMinutes = Math.ceil(recommendedMinutes * 1.1);
  }

  return {
    min: minMinutes,
    max: maxMinutes,
    recommended: recommendedMinutes,
    guideline,
  };
}

/**
 * 获取喂奶建议文本
 * @param ageInMonths 宝宝月龄
 * @returns 为父母准备的建议文本
 */
export function getFeedingRecommendationText(ageInMonths: number): string {
  const guideline = getFeedingGuidelineByAge(ageInMonths);

  const feedingTypeText = {
    demand: '按需喂养',
    scheduled: '定时喂养',
    mixed: '按需与定时相结合',
  };

  const nightFeedingText = {
    continue: '继续夜间喂奶',
    reduce: '逐步减少夜间喂奶',
    stop: '停止夜间喂奶',
  };

  return `
建议喂奶方式: ${feedingTypeText[guideline.feedingType]}
建议日喂奶次数: ${guideline.recommendedTimesPerDay} 次
建议间隔: ${guideline.intervalMinHours}-${guideline.intervalMaxHours} 小时
${nightFeedingText[guideline.nightFeeding]}
  `.trim();
}

/**
 * 计算宝宝月龄
 * @param birthDate 出生日期（YYYY-MM-DD 格式或 Date 对象）
 * @returns 精确的月龄（包括小数）
 */
export function calculateAgeInMonths(birthDate: string | Date): number {
  const birth = typeof birthDate === 'string'
    ? new Date(birthDate)
    : birthDate;

  const now = new Date();

  // 计算完整月数
  let months = (now.getFullYear() - birth.getFullYear()) * 12;
  months += now.getMonth() - birth.getMonth();

  // 如果还未达到相同日期，减去一个月
  if (now.getDate() < birth.getDate()) {
    months--;
  }

  // 计算天数用于小数部分
  const lastMonthDate = new Date(now.getFullYear(), now.getMonth(), birth.getDate());
  const daysIntoMonth = (now.getTime() - lastMonthDate.getTime()) / (1000 * 60 * 60 * 24);
  const daysInMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate();
  const monthFraction = daysIntoMonth / daysInMonth;

  return Math.max(0, months + monthFraction);
}
