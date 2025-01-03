import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import Cookies from 'js-cookie';

import './assets/base.css';
import './assets/main.css'; 

import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';

const savedTheme = Cookies.get('theme');
if (savedTheme === 'dark') {
  document.body.classList.add('dark-theme');
}

const app = createApp(App);
app.use(router).mount('#app');
