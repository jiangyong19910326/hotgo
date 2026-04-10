<template>
  <n-config-provider :theme="null" :theme-overrides="themeOverrides" :locale="naiveLocale" :date-locale="naiveDateLocale">
    <n-message-provider placement="top-right" :max="4">
      <n-notification-provider placement="top-right">
        <router-view />
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  NConfigProvider, NMessageProvider, NNotificationProvider,
  zhCN, dateZhCN, enUS, dateEnUS,
} from 'naive-ui';
import type { GlobalThemeOverrides } from 'naive-ui';

const { locale } = useI18n();
const naiveLocale     = computed(() => locale.value === 'en-US' ? enUS     : zhCN);
const naiveDateLocale = computed(() => locale.value === 'en-US' ? dateEnUS : dateZhCN);

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor:         '#e05535',
    primaryColorHover:    '#f0795a',
    primaryColorPressed:  '#c43e22',
    primaryColorSuppl:    '#e05535',
    bodyColor:            '#f5f0e8',
    cardColor:            '#ffffff',
    modalColor:           '#ffffff',
    popoverColor:         '#ffffff',
    tableColor:           '#ffffff',
    tableColorHover:      '#faf5ee',
    tableHeaderColor:     '#f0ebe0',
    borderColor:          'rgba(28, 22, 13, 0.18)',
    borderRadius:         '0px',
    borderRadiusSmall:    '0px',
    textColorBase:        '#1c160d',
    textColor1:           '#1c160d',
    textColor2:           '#6b5a4e',
    textColor3:           '#b5a898',
    placeholderColor:     '#b5a898',
    fontFamily:           '"DM Sans", system-ui, sans-serif',
    fontFamilyMono:       '"JetBrains Mono", monospace',
    fontWeight:           '400',
    fontWeightStrong:     '600',
    lineHeight:           '1.6',
    heightSmall:          '30px',
    heightMedium:         '38px',
    heightLarge:          '46px',
  },
  Input: {
    color:              '#ffffff',
    colorFocus:         '#ffffff',
    border:             '1px solid #1c160d',
    borderHover:        '1px solid #1c160d',
    borderFocus:        '1px solid #e05535',
    boxShadowFocus:     '3px 3px 0 rgba(224, 85, 53, 0.12)',
    textColor:          '#1c160d',
    placeholderColor:   '#b5a898',
    caretColor:         '#e05535',
    clearColor:         '#b5a898',
    borderRadius:       '0px',
    fontSizeLarge:      '15px',
  },
  Button: {
    textColor:           '#6b5a4e',
    textColorHover:      '#1c160d',
    textColorPressed:    '#1c160d',
    borderColor:         'rgba(28, 22, 13, 0.3)',
    borderColorHover:    '#1c160d',
    colorHover:          '#f0ebe0',
    borderRadius:        '0px',
    fontWeight:          '500',
    paddingSmall:        '0 14px',
    paddingMedium:       '0 20px',
    paddingLarge:        '0 28px',
  },
  DataTable: {
    thColor:         '#f0ebe0',
    thTextColor:     '#6b5a4e',
    tdColor:         '#ffffff',
    tdColorHover:    '#faf5ee',
    tdTextColor:     '#1c160d',
    borderColor:     'rgba(28, 22, 13, 0.1)',
    thFontWeight:    '600',
    fontSizeSmall:   '13px',
    fontSizeMedium:  '13.5px',
    borderRadius:    '0px',
  },
  Form: {
    labelTextColor:         '#6b5a4e',
    feedbackTextColorError: '#c0392b',
  },
  DatePicker: {
    panelColor:           '#ffffff',
    itemColorActive:      '#e05535',
    itemColorHover:       '#faf5ee',
    itemTextColorActive:  '#ffffff',
    borderRadius:         '0px',
  },
  Message: {
    colorInfo:         '#ffffff',
    colorSuccess:      '#f0faf5',
    colorError:        '#fff5f3',
    colorWarning:      '#fffbf0',
    textColorInfo:     '#1c160d',
    textColorSuccess:  '#2d8a5e',
    textColorError:    '#c0392b',
    textColorWarning:  '#b87a16',
    borderInfo:        '1px solid rgba(28,22,13,0.15)',
    borderSuccess:     '1px solid rgba(45,138,94,0.25)',
    borderError:       '1px solid rgba(192,57,43,0.25)',
    borderWarning:     '1px solid rgba(184,122,22,0.25)',
    borderRadius:      '0px',
    boxShadow:         '4px 4px 0 rgba(28, 22, 13, 0.08)',
  },
  Tag: {
    colorSuccess:      '#f0faf5',
    textColorSuccess:  '#2d8a5e',
    borderSuccess:     '1px solid rgba(45,138,94,0.3)',
    colorError:        '#fff5f3',
    textColorError:    '#c0392b',
    borderError:       '1px solid rgba(192,57,43,0.25)',
    colorWarning:      '#fffbf0',
    textColorWarning:  '#b87a16',
    borderWarning:     '1px solid rgba(184,122,22,0.25)',
    colorInfo:         '#fff8f5',
    textColorInfo:     '#e05535',
    borderInfo:        '1px solid rgba(224,85,53,0.25)',
    fontFamily:        '"JetBrains Mono", monospace',
    fontSizeSmall:     '11px',
    borderRadius:      '0px',
    padding:           '0 7px',
  },
  Progress: {
    railColor:        'rgba(28, 22, 13, 0.08)',
    fillColor:        '#e05535',
    fillColorInfo:    '#e05535',
    fillColorSuccess: '#2d8a5e',
    borderRadius:     '0px',
  },
  Spin: { color: '#e05535' },
  Pagination: {
    itemColorHover:      '#faf5ee',
    itemColorActive:     '#1c160d',
    itemTextColorActive: '#f5f0e8',
    itemBorderActive:    '1px solid #1c160d',
    buttonColor:         '#ffffff',
    buttonColorHover:    '#faf5ee',
    buttonBorder:        '1px solid rgba(28,22,13,0.2)',
    buttonBorderHover:   '1px solid #1c160d',
    borderRadius:        '0px',
  },
};
</script>
