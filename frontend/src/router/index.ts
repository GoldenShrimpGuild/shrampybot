import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useGlobalStore } from '../stores/global-store'
import { useAuthStore } from '../stores/auth'
import { useUserStore } from '../stores/user'
import AuthLayout from '../layouts/AuthLayout.vue'
import AppLayout from '../layouts/AppLayout.vue'
import axios from 'axios'

import RouteViewComponent from '../layouts/RouterBypass.vue'

export interface INavigationRoute {
  name?: string
  path: string
  meta: {
    nav: {
      icon?: string
      displayName: string
      disabled: boolean
      hidden: boolean
    }
    perms: {
      requiresAuth?: boolean
      requiresStaff?: boolean
      requiresAdmin?: boolean
    }
  }
  redirect?: object
  component?: object
  children?: INavigationRoute[]
}

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    redirect: { name: 'streams' },
    meta: {
      nav: {
        icon: 'home',
        displayName: 'menu.home',
        disabled: false,
        hidden: false,
      },
      perms: {
        requiresAuth: false,
        requiresStaff: false,
        requiresAdmin: false,
      },
    },
    component: AppLayout,
    children: [
      {
        name: 'streams',
        path: '/streams',
        meta: {
          nav: {
            icon: '',
            displayName: 'menu.activeStreams',
            disabled: false,
            hidden: false,
          },
          perms: {
            requiresAuth: true,
            requiresStaff: false,
            requiresAdmin: false,
          },
        },
        component: () => import('../pages/public/ActiveStreams.vue'),
      },
      // {
      //   name: 'dashboard',
      //   path: '/dashboard',
      //   meta: {
      //     nav: {
      //       icon: '',
      //       displayName: 'menu.dashboard',
      //       disabled: false,
      //       hidden: false,
      //     },
      //     perms: {
      //       requiresAuth: true,
      //       requiresStaff: false,
      //       requiresAdmin: false,
      //     },
      //   },
      //   component: () => import('../pages/admin/dashboard/Dashboard.vue'),
      // },
    ],
  },
  {
    path: '/admin',
    redirect: { name: 'users' },
    meta: {
      nav: {
        icon: 'lock',
        displayName: 'menu.admin',
        disabled: false,
        hidden: false,
      },
      perms: {
        requiresAuth: false,
        requiresStaff: false,
        requiresAdmin: true,
      },
    },
    component: AppLayout,
    children: [
      {
        name: 'user_list',
        path: '/admin/users',
        meta: {
          nav: {
            icon: '',
            displayName: 'menu.user_list',
            disabled: false,
            hidden: false,
          },
          perms: {
            requiresAuth: true,
            requiresStaff: false,
            requiresAdmin: true,
          },
        },
        component: () => import('../pages/admin/UserList.vue'),
      },
    ],
  },
  {
    path: '/auth',
    meta: {
      nav: {
        icon: 'vuestic-iconset-dashboard',
        displayName: 'menu.auth',
        disabled: true,
        hidden: true,
      },
      perms: {
        requiresAuth: false,
        requiresStaff: false,
        requiresAdmin: false,
      },
    },
    component: AuthLayout,
    children: [
      {
        name: 'login',
        path: 'login',
        meta: {
          nav: {
            icon: 'vuestic-iconset-dashboard',
            displayName: 'menu.login',
            disabled: true,
            hidden: false,
          },
          perms: {
            requiresAuth: false,
            requiresStaff: false,
            requiresAdmin: false,
          },
        },
        component: () => import('../pages/auth/Login.vue'),
      },
      {
        name: 'logout',
        path: 'logout',
        meta: {
          nav: {
            icon: 'vuestic-iconset-dashboard',
            displayName: 'menu.logout',
            disabled: true,
            hidden: true,
          },
          perms: {
            requiresAuth: true,
            requiresStaff: false,
            requiresAdmin: false,
          },
        },
        component: () => import('../pages/auth/Logout.vue'),
      },
      {
        name: 'validate_oauth',
        path: 'validate_oauth',
        meta: {
          nav: {
            icon: 'vuestic-iconset-dashboard',
            displayName: 'menu.activeStreams',
            disabled: true,
            hidden: true,
          },
          perms: {
            requiresAuth: false,
            requiresStaff: false,
            requiresAdmin: false,
          },
        },
        component: () => import('../pages/auth/ValidateOAuth.vue'),
      },
      {
        path: '',
        meta: {
          nav: {
            icon: 'vuestic-iconset-dashboard',
            displayName: 'menu.activeStreams',
            disabled: true,
            hidden: true,
          },
          perms: {
            requiresAuth: false,
            requiresStaff: false,
            requiresAdmin: false,
          },
        },
        redirect: { name: 'login' },
      },
    ],
  },
  {
    name: '404',
    path: '/404',
    component: () => import('../pages/404.vue'),
    meta: {
      nav: {
        icon: 'vuestic-iconset-dashboard',
        displayName: 'menu.activeStreams',
        disabled: true,
        hidden: true,
      },
      perms: {
        requiresAuth: false,
        requiresStaff: false,
        requiresAdmin: false,
      },
    },
  },
]

export const navRoutes = routes as INavigationRoute[]

const router = createRouter({
  history: createWebHistory('/shrampybot/'),
  // history: createWebHashHistory('/shrampybot/'),
  // scrollBehavior(to, from, savedPosition) {
  //   if (savedPosition) {
  //     return savedPosition
  //   }
  //   // For some reason using documentation example doesn't scroll on page navigation.
  //   if (to.hash) {
  //     return { el: to.hash, behavior: 'smooth' }
  //   } else {
  //     window.scrollTo(0, 0)
  //   }
  // },
  routes,
})

router.beforeEach(async (to: any, from: any, next) => {
  const AuthStore = useAuthStore()
  const UserStore = useUserStore()

  if (AuthStore.accessToken !== '') {
    next()
    return
  }

  // instead of having to check every route record with
  // to.matched.some(record => record.meta.requiresAuth)
  if (to.meta.perms.requiresAuth) {
    if (AuthStore.accessToken !== '') {
      next()
      // next(await validateAndFetchRoute(to))
    } else {
      // Try refreshing if no accessToken is set.
      await AuthStore.callRefresh()

      // If there's still no accessToken set after calling refresh, route to auth screen
      if (AuthStore.accessToken === '') {
        next('/auth/login')
      } else {
        next()
      }
    }
    return
    // this route requires auth, check if logged in
    // if not, redirect to login page.
  }
  next()
})

export const validateAndFetchRoute = async (route_path: any) => {
  const AuthStore = useAuthStore()
  const UserStore = useUserStore()
  const GlobalStore = useGlobalStore()

  const path = '/auth/self'

  const axiosConfig = AuthStore.getAxiosConfig()

  try {
    const bearerResponse = await axios.get(path, axiosConfig)
    UserStore.$state.self = bearerResponse.data.self
  } catch (error: any) {
    if (error.response.status in [400, 401]) {
      const refresh_path = '/token/refresh/'

      try {
        const refreshResponse = await axios.post(
          refresh_path,
          {},
          {
            baseURL: GlobalStore.getApiBaseUrl(),
            withCredentials: true,
            headers: {
              'Content-Type': 'application/json',
            },
          },
        )
        AuthStore.$state.accessToken = refreshResponse.data.access
        route_path = await validateAndFetchRoute(route_path)
      } catch (refreshError: any) {
        AuthStore.$state.accessToken = ''
        route_path = {
          name: 'login',
          path: '/auth/login',
        } as RouteRecordRaw
      }
    } else if (error.response.status == 500) {
      AuthStore.$state.accessToken = ''
      route_path = {
        name: 'login',
        path: '/auth/login',
      } as RouteRecordRaw
    }
  }
  return route_path
}

export default router
