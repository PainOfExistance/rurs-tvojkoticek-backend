<script setup>
import Home from "@/components/Home.vue";
import Navbar from "@/components/Navbar.vue";
import Header from "@/components/Header.vue";
import Footer from "@/components/Footer.vue";
</script>

<template>
  <Header/>
  <Navbar/>
  <div class="container">
    <router-view style="position: relative"/>
  </div>
  <Footer/>
</template>

<style scoped>
header {
  line-height: 1.5;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
</style>
