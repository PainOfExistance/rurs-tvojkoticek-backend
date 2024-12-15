<template>
  <div>
    <button v-if="!menuOpen" @click="toggleMenu" class="menu-button">
      <span>
        ⦿
        <br>
        ⦿
        <br>
        ⦿
      </span>
    </button>

    <transition name="slide">
      <div v-if="menuOpen" class="side-menu">
        <ul>
          <li @click="toggleMenu"><span>Zapri meni</span></li>
          <hr>
          <li><router-link to="/">Domov</router-link></li>
          <li><router-link to="/posts">Obstoječa vprašanja</router-link></li>
          <li><router-link to="/ask">Postavi vprašanje</router-link></li>
          <li v-if="!isLoggedIn"><router-link to="/login">Prijava</router-link></li>
          <li v-if="!isLoggedIn"><router-link to="/register">Registracija</router-link></li>
          <li v-if="isLoggedIn"><router-link to="/profile">Profil</router-link></li>
          <hr v-if="isLoggedIn">
          <li v-if="isLoggedIn"><router-link to="/" @click="logout">Odjava</router-link></li>
        </ul>
      </div>
    </transition>
  </div>
</template>

<script>
import Cookies from 'js-cookie';
export default {
  data() {
    return {
      menuOpen: false,
      isLoggedIn: false,
    };
  },
  methods: {
    checkLoginStatus() {
      this.isLoggedIn = !!Cookies.get('id');
    },
    toggleMenu() {
      this.menuOpen = !this.menuOpen;
    },
    logout() {
      Cookies.remove('id', { path: '/' });
      this.isLoggedIn = false;
    }
  },
  mounted() {
    this.checkLoginStatus();
  },
};
</script>

<style scoped>
.menu-button {
  position: fixed;
  left: 10px;
  top: 45%;
  font-size: 20px;
  padding: 8px 8px;
  cursor: pointer;
  z-index: 1001;
  background-color: #54627b;
  color: white;
  border: none;
  border-radius: 5px;
}

.side-menu {
  position: fixed;
  top: 35%;
  left: 10px;
  width: 200px;
  background-color: #54627b;
  color: #fff;
  overflow-x: hidden;
  padding: 10px 20px;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.3);
  z-index: 1000;
  border-radius: 15px;
}

.side-menu ul {
  list-style-type: none;
  padding: 0;
}

.side-menu li {
  margin: 10px 0;
}

.side-menu span:hover {
  cursor: pointer;
}

.side-menu a,
.side-menu span {
  color: #fff;
  text-decoration: none;
  font-size: 18px;
}
</style>
