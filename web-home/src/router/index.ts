import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/views/layout/MainLayout.vue'),
      children: [
        {
          path: '',
          name: 'Home',
          component: () => import('@/views/Home.vue'),
          meta: { title: '短链首页' },
        },
        {
          path: 'list',
          name: 'List',
          component: () => import('@/views/List.vue'),
          meta: { title: '短链列表' },
        },
        {
          path: 'stats/:code',
          name: 'Stats',
          component: () => import('@/views/Stats.vue'),
          meta: { title: '访问统计' },
        },
      ],
    },
    {
      // 短链跳转：访问 /:code 时由此路由处理
      path: '/r/:code',
      name: 'Redirect',
      component: () => import('@/views/Redirect.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
});

router.afterEach((to) => {
  const title = (to.meta?.title as string) || '短链服务';
  document.title = `${title} - 短链服务`;
});

export default router;
