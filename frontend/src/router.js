import Answer from "@/components/Answer/Answer.vue";
import Answers from "@/components/Answer/Answers.vue";
import Posts from "@/components/Posts/Posts.vue";
import Ask from "@/components/Question/Ask.vue";
import Details from "@/components/Question/Details.vue";
import Profile from "@/components/user/Profile.vue";
import PostVideo from "@/components/Videostore/PostVideo.vue";
import { createRouter, createWebHistory } from 'vue-router';
import Home from './components/Home.vue';
import Landing from './components/Landing.vue';
import Login from './components/user/Login.vue';
import Register from './components/user/Register.vue';
import Videostore from './components/Videostore/VideoPage.vue';

import Cookies from 'js-cookie';

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: Landing,
    meta: { hideHeader: true, requiresGuest: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { hideHeader: true, requiresGuest: true },
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { hideHeader: true, requiresGuest: true },
  },
  {
    path: '/home',
    name: 'Home',
    component: Home,
    meta: { hideNavbar: true, hideFooter: true, requiresAuth: true } // Prikazuj Header, skrži Navbar in Footer, zahteva prijavo
  },
  {
    path: '/ask',
    name: 'Ask',
    component: Ask,
    meta: { requiresAuth: true }, // Zahteva prijavo
  },
  {
    path: '/post-video',
    name: 'PostVideo',
    component: PostVideo,
    meta: { requiresAuth: true }, // Zahteva prijavo
  },
  {
    path: '/posts',
    name: 'Posts',
    component: Posts,
  },
  {
    path: '/videostore',
    name: 'Videostore',
    component: Videostore,
  },
  {
    path: '/details/:id',
    name: 'Details',
    component: Details,
    props: true,
  },
  {
    path: '/answers',
    name: 'Answers',
    component: Answers,
  },
  {
    path: '/answer',
    name: 'Answer',
    component: Answer,
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }, 
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('./components/user/Settings.vue'),
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const isLoggedIn = !!Cookies.get('id');

  if (to.meta.requiresAuth && !isLoggedIn) {
    next({ name: 'Landing' });
  } else if (to.meta.requiresGuest && isLoggedIn) {
    next({ name: 'Home' });
  } else {
    next();
  }
});

export default router;
