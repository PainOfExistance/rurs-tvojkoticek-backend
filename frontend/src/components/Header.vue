<template>
  <div class="header">
    <!-- Logo Section -->
    <div class="logo-section">
      <router-link to="/home">
        <img class="logo-image" alt="Logo" src="@/assets/imgs/logo.png" />
      </router-link>
    </div>

    <!-- Navigation Links -->
    <nav class="navigation">
      <router-link to="/home" class="nav-link">Domača Stran</router-link>
      <router-link to="/ask" class="nav-link">Vprašaj</router-link>
      <router-link to="/posts" class="nav-link">Objave</router-link>
      <router-link to="/post-video" class="nav-link">Objavi Video</router-link>
      <router-link to="/videostore" class="nav-link">Videoteka</router-link>
    </nav>

    <!-- Profile Dropdown -->
    <div class="profile-section dropdown">
      <!-- User Icon for Dropdown -->
      <button
        class="btn dropdown-toggle p-0 border-0"
        type="button"
        id="profileDropdown"
        data-bs-toggle="dropdown"
        aria-expanded="false"
      >
        <img 
          class="profile-image" 
          :src="profileImage" 
          @error="setDefaultImage"
          alt="User Icon"
        />
      </button>

      <!-- Dropdown Menu -->
      <ul class="dropdown-menu dropdown-menu-end" aria-labelledby="profileDropdown">
        <li>
          <router-link to="/profile" class="dropdown-item">
            <i class="bi bi-person-fill me-2"></i> Profil
          </router-link>
        </li>
        <li>
          <router-link to="/settings" class="dropdown-item">
            <i class="bi bi-gear-fill me-2"></i> Nastavitve
          </router-link>
        </li>
        <li><hr class="dropdown-divider" /></li>
        <li>
          <button @click="logout" class="dropdown-item">
            <i class="bi bi-box-arrow-right me-2"></i> Odjava
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import Cookies from 'js-cookie';

export default {
  name: "Header",
  data() {
    return {
      profileImage: "@/assets/imgs/profile_icon.png", // Default image
    };
  },
  methods: {
    logout() {
      // Remove cookies
      Cookies.remove('id');
      Cookies.remove('username');
      // Redirect to landing page
      this.$router.push({ name: 'Landing' });
    },
    setDefaultImage(event) {
      // Set placeholder image if the current one can't load
      event.target.src = "https://via.placeholder.com/45?text=User";
    },
  },
};
</script>

<style scoped>
@import url('https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.5/font/bootstrap-icons.css');

.header {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 80px;
  background-color: rgba(255, 255, 255, 0.95);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  font-family: 'Poppins', sans-serif;
}

.logo-section {
  display: flex;
  align-items: center;
}

.logo-image {
  height: 60px;
  width: auto;
}

.navigation {
  display: flex;
  gap: 20px;
}

.nav-link {
  text-decoration: none;
  color: #444;
  font-size: 1rem;
  font-weight: 500;
  transition: color 0.3s ease;
}

.nav-link:hover {
  color: #8da4b9;
}

.profile-section {
  display: flex;
  align-items: center;
}

.profile-image {
  height: 45px;
  width: 45px;
  border-radius: 50%;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.profile-image:hover {
  transform: scale(1.05);
}

.dropdown-menu {
  min-width: 200px;
}

.dropdown-item {
  display: flex;
  align-items: center;
  font-size: 0.95rem;
  color: #444;
}

.dropdown-item:hover {
  background-color: #f8f9fa;
}

.dropdown-item i {
  font-size: 1.1rem;
}
</style>
