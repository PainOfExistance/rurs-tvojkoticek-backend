<template>
  <div>
      <div class="top-menu">
        <img class="logo-image" alt="No logo :c" src="../assets/imgs/logo.png">
        <router-link class="profile-link" to="/profile"><img class="profile-image" src="../assets/profile_icon.png" alt="No profile :c"></router-link>
      </div>
  </div>
</template>

<script>
export default {
  data() {
    return {

    };
  },
  methods: {

  },
};
</script>

<style scoped>
.top-menu {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 150px;
  padding: 20px;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.3);
  background-color: white;
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.logo-image{
  height: 150%;
  max-width: 100%;
}

.profile-link {
  position: fixed;
  right: 1%;
}

.profile-image {
  height: 45px;
  width: 45px;
}
</style>
