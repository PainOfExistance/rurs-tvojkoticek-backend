<template>
  <div class="settings-page">
    <h2 class="title">Nastavitve</h2>
    <div class="toggle-container">
      <label class="theme-toggle">
        <!-- Proper input toggle -->
        <input 
          type="checkbox" 
          v-model="isDarkTheme" 
          @change="toggleTheme" 
          aria-label="Toggle Dark Theme"
        />
        <span class="slider round"></span>
      </label>
      <span class="toggle-label">{{ isDarkTheme ? 'Dark Theme' : 'Light Theme' }}</span>
    </div>
  </div>
</template>

<script>
import Cookies from 'js-cookie';

export default {
  name: "Settings",
  data() {
    return {
      isDarkTheme: false, // Default to light theme
    };
  },
  methods: {
    toggleTheme() {
      // Apply or remove the 'dark-theme' class on the body
      if (this.isDarkTheme) {
        document.body.classList.add('dark-theme');
      } else {
        document.body.classList.remove('dark-theme');
      }
      // Save the theme preference in cookies
      Cookies.set('theme', this.isDarkTheme ? 'dark' : 'light');
    },
    applyInitialTheme() {
      if (this.isDarkTheme) {
        document.body.classList.add('dark-theme');
      } else {
        document.body.classList.remove('dark-theme');
      }
    }
  },
  mounted() {
    this.isDarkTheme = Cookies.get('theme') === 'dark';
    this.applyInitialTheme();
  },
};
</script>

<!-- Remove 'scoped' to allow global styles -->
<style>
.settings-page {
  text-align: center;
  padding: 2rem;
  transition: background-color 0.3s ease, color 0.3s ease;
}

/* Global Theme-specific styles */
body.light-theme {
  background-color: #f4f4f4;
  color: #333333;
}

body.dark-theme {
  background-color: #121212;
  color: #f1f1f1;
}

.title {
  font-size: 2rem;
  font-weight: 600;
}

.toggle-container {
  margin-top: 2rem;
  display: flex;
  align-items: center;
  gap: 10px;
}

.theme-toggle {
  position: relative;
  display: inline-block;
  width: 60px; /* Increased size for easier clicking */
  height: 34px;
}

.theme-toggle input {
  display: none; /* Using display none to remove input but keep accessible for screen readers */
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
  border-radius: 34px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #8da4b9;
}

input:checked + .slider:before {
  transform: translateX(26px);
}

.toggle-label {
  font-size: 1.2rem;
}
</style>
