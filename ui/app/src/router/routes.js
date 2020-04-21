import MainLayout from 'layouts/MainLayout.vue'
import Login from 'pages/Login.vue'
import DesktopTemplates from 'pages/DesktopTemplates.vue'
import VNCViewer from 'pages/VNCViewer.vue'
import Settings from 'pages/Settings.vue'
import Profile from 'pages/Profile.vue'
import Error404 from 'pages/Error404.vue'

const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      // {
      //   path: '',
      //   name: '',
      //   component: DesktopTemplates,
      //   meta: { requiresAuth: true }
      // },
      {
        path: 'login',
        name: 'login',
        component: Login
      },
      {
        path: 'templates',
        name: 'templates',
        component: DesktopTemplates,
        meta: { requiresAuth: true }
      },
      {
        path: 'control',
        name: 'control',
        component: VNCViewer,
        meta: { requiresAuth: true }
      },
      {
        path: 'settings',
        name: 'settings',
        component: Settings,
        meta: { requiresAuth: true }
      },
      {
        path: 'profile',
        name: 'profile',
        component: Profile,
        meta: { requiresAuth: true }
      },
      { path: '*', component: Error404 }
    ]
  }
]

// Always leave this as last one
if (process.env.MODE !== 'ssr') {
  routes.push({
    path: '*',
    component: () => import('pages/Error404.vue')
  })
}

export default routes
