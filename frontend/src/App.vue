<!-- src/App.vue -->
<template>
  <div id="app">
    <!-- Example header or navigation -->
    <Header v-if="!$route.meta.hideHeader" />
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
    $route(to) {
      // If route meta says hideHeader => prevent scrolling
      if (to.meta.hideHeader) {
        document.body.classList.add('no-scroll');
      } else {
        document.body.classList.remove('no-scroll');
      }
    },
  },
  mounted() {
    // Initial scrolling logic
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
/* Minimal styling for layout */
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* Provide some padding if not hiding header */
.container {
  flex-grow: 1;
  padding: 20px;
}

.no-container {
  padding: 0;
}

/* Prevent scrolling if 'no-scroll' is active on body */
.no-scroll {
  overflow: hidden;
  height: 100%;
}
</style>
