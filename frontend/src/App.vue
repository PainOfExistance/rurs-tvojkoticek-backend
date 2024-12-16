<template>
  <div id="app">
    <!-- Show Header only if not hidden via route meta -->
    <Header v-if="!$route.meta.hideHeader" />

    <!-- Main Content -->
    <div :class="containerClass">
      <router-view />
    </div>
  </div>
</template>

<script>
import Header from './components/Header.vue';

export default {
  name: 'App',
  computed: {
    containerClass() {
      return this.$route.meta.hideHeader ? 'no-container' : 'container';
    },
  },
  watch: {
    // Watch route changes to toggle the 'no-scroll' class on body
    $route(to) {
      if (to.meta.hideHeader) {
        document.body.classList.add('no-scroll');
      } else {
        document.body.classList.remove('no-scroll');
      }
    },
  },
  mounted() {
    // Set correct 'no-scroll' class on initial load
    if (this.$route.meta.hideHeader) {
      document.body.classList.add('no-scroll');
    } else {
      document.body.classList.remove('no-scroll');
    }
  },
  components: {
    Header,
  },
};
</script>

<style>
/* Box-sizing and reset styles */
*, *::before, *::after {
  box-sizing: border-box;
}

body, html {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: 'Poppins', sans-serif;
  background: linear-gradient(135deg, #fbcfe8, #fde2e4, #fbd4e7, #fbe9e7);
  background-size: 400% 400%;
  animation: gradientBG 20s ease infinite;
}


#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: transparent;
}

.container {
  flex-grow: 1;
  padding: 20px; /* Default padding */
}

.no-container {
  padding: 0; /* Remove padding on pages without header */
}

/* Disable scrolling */
.no-scroll {
  overflow: hidden;
  height: 100%;
}

@keyframes gradientBG {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}
</style>
