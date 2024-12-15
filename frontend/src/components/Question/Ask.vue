<!-- src/components/Question/Ask.vue -->
<template>
  <div class="home-page">
    <!-- Dekorativne Oblike -->
    <div class="shape shape1"></div>
    <div class="shape shape2"></div>
    <div class="shape shape3"></div>

    <div class="content-wrapper">
      <img class="app-logo" src="@/assets/imgs/logo.png" alt="Logo" />

      <h2 class="form-title">Postavite svoje vprašanje</h2>

      <div v-if="isLoggedIn" class="post-add-container">
        <form @submit.prevent="submitAQuestion" class="ask-form">
          <div class="form-group">
            <label for="authorInput" class="form-label">Avtor:</label>
            <input
              id="authorInput"
              class="form-input"
              placeholder="Vnesite vaše ime..."
              required
              v-model="formData.name"
            />
          </div>

          <div class="form-group">
            <label for="storyInput" class="form-label">Vaše vprašanje:</label>
            <textarea
              id="storyInput"
              class="form-textarea"
              rows="5"
              placeholder="Opišite vaš problem in kaj vas zanima..."
              required
              v-model="formData.question"
            ></textarea>
          </div>

          <button type="submit" class="btn btn-primary">Vprašaj</button>
        </form>
      </div>

      <div v-else class="not-logged-in-section">
        <p class="not-logged-in-info">
          V kolikor želite postaviti vprašanje, morate biti prijavljeni v sistem.
        </p>
        <div class="button-group">
          <router-link to="/login" class="btn btn-primary">
            Prijava
          </router-link>
          <router-link to="/register" class="btn btn-primary">
            Registracija
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Cookies from 'js-cookie';
import axios from 'axios';

export default {
  name: "Ask",
  data() {
    return {
      isLoggedIn: false,
      formData: {
        name: '',
        question: '',
      },
    };
  },
  methods: {
    checkLoginStatus() {
      this.isLoggedIn = !!Cookies.get('id');
    },
    async submitAQuestion() {
      try {
        const response = await axios.post('http://localhost:8080/post', {
          username: this.formData.name,
          problem: this.formData.question,
          date: new Date().toISOString().split('T')[0],
        });

        if (response.status === 200) {
          this.$router.push({ name: 'Posts' });
        }
      } catch (error) {
        console.error(error);
        // Dodajte obravnavo napak po potrebi
      }
    }
  },
  mounted() {
    this.checkLoginStatus();
  },
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap');

.home-page {
  position: relative;
  width: 100%;
  min-height: calc(100vh - 80px);
  background: linear-gradient(135deg);
  background-size: 400% 400%;
  animation: gradientBG 20s ease infinite;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: 'Poppins', sans-serif;
  color: #444;
  text-align: center;
  padding: 20px;
  overflow-y: auto;
  box-sizing: border-box;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.3;
  animation: float 25s infinite ease-in-out;
}

.shape1 {
  width: 200px;
  height: 200px;
  background: #fbd4e7;
  top: 10px;
  left: 10px;
}

.shape2 {
  width: 300px;
  height: 300px;
  background: #fde2e4;
  bottom: 10px;
  right: 10px;
}

.shape3 {
  width: 150px;
  height: 150px;
  background: #fbe9e7;
  top: 20%;
  right: 20%;
}

.content-wrapper {
  position: relative;
  z-index: 1;
  max-width: 800px;
  width: 90%;
  padding: 40px 30px;
  background: rgba(255, 255, 255, 0.25);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  animation: fadeIn 2s ease forwards;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 40px;
}

.app-logo {
  width: 150px;
  margin-bottom: 20px;
  animation: slideIn 1.5s ease forwards;
}

.form-title {
  font-size: 2rem;
  font-weight: 600;
  margin-bottom: 30px;
  color: #333;
  animation: fadeInUp 2s ease forwards;
}

.form-group {
  width: 100%;
  margin-bottom: 20px;
  text-align: left;
}

.form-label {
  font-size: 17px;
  font-weight: bold;
  margin-bottom: 8px;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 12px 15px;
  border: 1px solid #ccc;
  border-radius: 10px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-input:focus,
.form-textarea:focus {
  border-color: #8da4b9;
  outline: none;
  box-shadow: 0 0 5px rgba(138, 164, 185, 0.5);
}

.btn-primary {
  background-color: #8da4b9;
  border: none;
  color: white;
  padding: 15px 30px;
  border-radius: 30px;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.3s ease;
  text-decoration: none;
  display: inline-block;
  margin: 0 10px;
}

.btn-primary:hover {
  background-color: #7b93a9;
  transform: translateY(-3px);
}

.not-logged-in-section {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.not-logged-in-info {
  font-size: 22px;
  font-weight: bold;
  margin-bottom: 20px;
}

.button-group {
  display: flex;
  gap: 20px;
}

@keyframes gradientBG {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(20px); }
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slideIn {
  from { opacity: 0; transform: translateX(-50px); }
  to { opacity: 1; transform: translateX(0); }
}

/* Responsive Adjustments */
@media (max-width: 768px) {
  .content-wrapper {
    padding: 30px 20px;
    max-width: 90%;
    height: auto;
  }

  .form-title {
    font-size: 1.75rem;
  }

  .btn-primary {
    width: 100%;
    margin-bottom: 15px;
  }

  .button-group {
    flex-direction: column;
    gap: 10px;
    width: 100%;
  }
}
</style>