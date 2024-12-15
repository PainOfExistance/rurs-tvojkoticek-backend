import {createRouter, createWebHistory} from 'vue-router';
import Login from './components/user/Login.vue';
import Register from './components/user/Register.vue';
import Home from './components/Home.vue';
import App from "@/App.vue";
import Posts from "@/components/Posts/Posts.vue";
import Ask from "@/components/Question/Ask.vue";
import Details from "@/components/Question/Details.vue";
import Answer from "@/components/Answer/Answer.vue";
import Answers from "@/components/Answer/Answers.vue";
import Profile from "@/components/user/Profile.vue";

const routes = [
    {
        path: '/ask',
        name: 'Ask',
        component: Ask,
    },
    {
        path: '/posts',
        name: 'Posts',
        component: Posts,
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
        path: '/login',
        name: 'Login',
        component: Login,
    },
    {
        path: '/register',
        name: 'Register',
        component: Register,
    },
    {
        path: '/',
        name: 'Home',
        component: Home,
    },
    {
        path: '/profile',
        name: 'Profile',
        component: Profile,
    },

];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
