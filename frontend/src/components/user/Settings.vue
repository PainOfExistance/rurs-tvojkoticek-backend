<template>
  <div class="settings-page">
    <h2 class="title">Nastavitve</h2>
    <div class="toggle-container">
      <label class="theme-toggle">
        <input 
          type="checkbox"
          v-model="isDarkTheme"
          @change="toggleTheme"
        />
        <span class="slider round"></span>
      </label>
      <span class="toggle-label">
        {{ isDarkTheme ? 'Dark Theme' : 'Light Theme' }}
      </span>
    </div>
  </div>
</template>

<script>
import Cookies from 'js-cookie';

export default {
  name: 'Settings',
  data() {
    return {
      isDarkTheme: false,
    };
  },
  methods: {
    toggleTheme() {
      if (this.isDarkTheme) {
        document.body.classList.add('dark-theme');
      } else {
        document.body.classList.remove('dark-theme');
      }
      Cookies.set('theme', this.isDarkTheme ? 'dark' : 'light', { expires: 365 });
    },
    applyInitialTheme() {
      if (this.isDarkTheme) {
        document.body.classList.add('dark-theme');
      } else {
        document.body.classList.remove('dark-theme');
      }
    },
  },
  mounted() {
    const savedTheme = Cookies.get('theme');
    this.isDarkTheme = savedTheme === 'dark';
    this.applyInitialTheme();
  },
};
</script>

<style scoped>
.settings-page {
  text-align: center;
  padding: 2rem;
  transition: background-color 0.3s ease, color 0.3s ease;
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
  width: 60px;
  height: 34px;
}

.theme-toggle input {
  display: none;
}

.slider {
  position: absolute;
  inset: 0;
  background-color: #ccc;
  border-radius: 34px;
  cursor: pointer;
  transition: 0.4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background-color: #fff;
  transition: 0.4s;
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
  color: var(--color-text);
}
</style>
