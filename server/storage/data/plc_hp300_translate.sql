-- ============================================================
-- HP300 自动接入 point 名称中文化
-- 用法: 自动接入后, 执行此脚本把英文 name 改成中文
-- 默认匹配 host='113200006045', 多设备时可改 WHERE
-- ============================================================

UPDATE hg_plc_point p
JOIN hg_plc_device d ON d.id = p.device_id
SET p.name = CASE p.field
  -- 报警位
  WHEN 'AL_BLOWER_STARTER'        THEN '鼓风机故障'
  WHEN 'AL_COOLER_STARTER'        THEN '冷风机故障'
  WHEN 'AL_CRUSHER'               THEN '破碎机故障'
  WHEN 'AL_CRUSHER_ILLEGAL_START' THEN '非法启动'
  WHEN 'AL_CRUSHER_OVERLOAD'      THEN '破碎机过载'
  WHEN 'AL_DISCHARGER_SIGNAL'     THEN '排料信号异常'
  WHEN 'AL_EMGENCY_STOP'          THEN '紧急停机'
  WHEN 'AL_HEATER_STARTER'        THEN '加热器故障'
  WHEN 'AL_HIGH_LUBE_PRESS'       THEN '润滑压力高'
  WHEN 'AL_HIGH_RETURN_TEMP'      THEN '回油温度高'
  WHEN 'AL_HYD_CLOGGED'           THEN '液压过滤器堵塞'
  WHEN 'AL_HYD_EXC_PRESS_ATT'     THEN '液压压力异常'
  WHEN 'AL_HYD_LS_LOW'            THEN '液压液位低'
  WHEN 'AL_HYD_PRESS_TO'          THEN '液压压力超时'
  WHEN 'AL_HYD_PUMP'              THEN '液压泵故障'
  WHEN 'AL_LOW_HYD_PRESS'         THEN '液压压力低'
  WHEN 'AL_LOW_LUBE_PRESS'        THEN '润滑压力低'
  WHEN 'AL_LOW_RETURN_TEMP'       THEN '回油温度低'
  WHEN 'AL_LUBE_CLOGGED'          THEN '润滑过滤器堵塞'
  WHEN 'AL_LUBE_LS_LOW'           THEN '润滑液位低'
  WHEN 'AL_LUBE_PRESS_FAULT'      THEN '润滑压力故障'
  WHEN 'AL_LUBE_PRESS_SENSOR'     THEN '润滑压力传感器故障'
  WHEN 'AL_LUBE_PUMP_STARTER'     THEN '润滑泵故障'
  WHEN 'AL_LUBE_RETURN_LS_LOW'    THEN '润滑回油液位低'
  WHEN 'AL_PS1B_PRESS_LOW'        THEN 'PS1B压力低'
  WHEN 'AL_PS1_PRESS_SENSOR'      THEN 'PS1压力传感器故障'
  WHEN 'AL_PS2B_PRESS_LOW'        THEN 'PS2B压力低'
  WHEN 'AL_PS2_PRESS_SENSOR'      THEN 'PS2压力传感器故障'
  WHEN 'AL_RETURN_TEMP_SENSOR'    THEN '回油温度传感器故障'
  WHEN 'AL_TANK_TEMP_SENSOR'      THEN '油箱温度传感器故障'

  -- 运行状态
  WHEN 'BLOWER_STARTER'           THEN '鼓风机启动'
  WHEN 'COOLER_MOTOR'             THEN '冷风电机'
  WHEN 'CRUSHER_MOTOR'            THEN '破碎机电机'
  WHEN 'DISCHARGER_CONVEYOR'      THEN '排料皮带'
  WHEN 'FEEDER_MOTER'             THEN '给料机电机'
  WHEN 'HYD_CLOGGED'              THEN '液压过滤器堵塞状态'
  WHEN 'HYD_LS'                   THEN '液压液位状态'
  WHEN 'HYD_PUMP_MOTOR'           THEN '液压泵电机'
  WHEN 'LUBE_CLOGGED'             THEN '润滑过滤器堵塞状态'
  WHEN 'LUBE_LS'                  THEN '润滑液位状态'
  WHEN 'LUBE_PUMP_MOTOR'          THEN '润滑泵电机'
  WHEN 'LUBE_RETURN_LS'           THEN '润滑回油液位状态'
  WHEN 'SOL1'                     THEN '电磁阀1'
  WHEN 'SOL2-3'                   THEN '电磁阀2-3'
  WHEN 'SOL4'                     THEN '电磁阀4'
  WHEN 'SOL5'                     THEN '电磁阀5'
  WHEN 'SOL6'                     THEN '电磁阀6'
  WHEN 'SOL7'                     THEN '电磁阀7'
  WHEN 'SOL8'                     THEN '电磁阀8'

  -- 模拟量
  WHEN 'CRUSHER_AMPERE'           THEN '破碎机电流'
  WHEN 'CRUSHER_SETTING'          THEN '破碎机设置'
  WHEN 'CRUSHER_TIME_H'           THEN '破碎机运行小时'
  WHEN 'CRUSHER_TIME_M'           THEN '破碎机运行分钟'
  WHEN 'LUBE_FILTER_TIME_H'       THEN '润滑滤芯运行小时'
  WHEN 'LUBE_FILTER_TIME_M'       THEN '润滑滤芯运行分钟'
  WHEN 'LUBE_PUMP_TIME_H'         THEN '润滑泵运行小时'
  WHEN 'LUBE_PUMP_TIME_M'         THEN '润滑泵运行分钟'
  WHEN 'LUBE_POST_FILLTER_PRESS'  THEN '润滑滤后压力'
  WHEN 'LUBE_PRE_FILLTER_PRESS'   THEN '润滑滤前压力'
  WHEN 'LUBE_RETURN_TEMP'         THEN '润滑回油温度'
  WHEN 'LUBE_TANK_TEMP'           THEN '润滑油箱温度'
  WHEN 'PS1_PRESS'                THEN 'PS1压力'
  WHEN 'PS2_PRESS'                THEN 'PS2压力'

  ELSE p.name
END,
-- 顺便补单位 (针对模拟量字段)
p.unit = CASE p.field
  WHEN 'CRUSHER_AMPERE'           THEN 'A'
  WHEN 'CRUSHER_TIME_H'           THEN 'h'
  WHEN 'CRUSHER_TIME_M'           THEN 'min'
  WHEN 'LUBE_FILTER_TIME_H'       THEN 'h'
  WHEN 'LUBE_FILTER_TIME_M'       THEN 'min'
  WHEN 'LUBE_PUMP_TIME_H'         THEN 'h'
  WHEN 'LUBE_PUMP_TIME_M'         THEN 'min'
  WHEN 'LUBE_POST_FILLTER_PRESS'  THEN 'bar'
  WHEN 'LUBE_PRE_FILLTER_PRESS'   THEN 'bar'
  WHEN 'LUBE_RETURN_TEMP'         THEN '℃'
  WHEN 'LUBE_TANK_TEMP'           THEN '℃'
  WHEN 'PS1_PRESS'                THEN 'bar'
  WHEN 'PS2_PRESS'                THEN 'bar'
  ELSE p.unit
END
WHERE d.host = '113200006045';
