<!-- src/components/user/Register.vue -->
<template>
    <div class="register-page">
      <!-- Dekorativne Oblike -->
      <div class="shape shape1"></div>
      <div class="shape shape2"></div>
      <div class="shape shape3"></div>
  
      <div class="content-wrapper">
        <img class="app-logo" src="@/assets/imgs/logo.png" alt="Logo" />
  
        <h2 class="form-title">Registracija</h2>
  
        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>
  
        <form @submit.prevent="addUser" class="register-form">
          <div class="form-group">
            <label for="username" class="form-label">Uporabniško ime:</label>
            <input
              type="text"
              id="username"
              class="form-input"
              required
              v-model="username"
              placeholder="Vnesi uporabniško ime..."
            />
          </div>
  
          <div class="form-group">
            <label for="email" class="form-label">Elektronski naslov:</label>
            <input
              type="email"
              id="email"
              class="form-input"
              required
              v-model="email"
              placeholder="Vnesi elektronski naslov..."
            />
          </div>
  
          <div class="form-group">
            <label for="password" class="form-label">Geslo:</label>
            <input
              type="password"
              id="password"
              class="form-input"
              required
              v-model="password"
              placeholder="Vnesi geslo..."
            />
          </div>
  
          <div class="form-group">
            <label for="repeat_password" class="form-label">Ponovi geslo:</label>
            <input
              type="password"
              id="repeat_password"
              class="form-input"
              required
              v-model="repeat_password"
              placeholder="Ponovno vnesi geslo..."
            />
          </div>
  
          <button type="submit" class="action-button register-button">Registracija</button>
        </form>
  
        <div class="additional-links">
          <p>Že imate račun? <router-link to="/login" class="link">Prijavite se</router-link></p>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import axios from "axios";
  import Cookies from "js-cookie";
  
  export default {
    name: "Register",
    data() {
      return {
        username: "",
        email: "",
        password: "",
        repeat_password: "",
        errorMessage: "",
      };
    },
    methods: {
      async addUser() {
        // Resetiraj sporočilo o napaki
        this.errorMessage = "";
  
        // Preveri, če se gesli ujemata
        if (this.password !== this.repeat_password) {
          this.errorMessage = "Gesli se ne ujemata";
          return;
        }
  
        // Preveri, če so vsa polja izpolnjena
        if (!this.username || !this.email || !this.password) {
          this.errorMessage = "Prosim, izpolnite vsa polja";
          return;
        }
  
        try {
          const response = await axios.post("http://localhost:8080/register", {
            username: this.username, 
            email: this.email,       
            password: this.password, 
          });
  
          // Preveri, če je registracija uspešna
          if (response.status === 201 || response.status === 200) {
            // Set cookies to mark user as logged in
            const expiryDate = new Date();
            expiryDate.setDate(expiryDate.getDate() + 7);
            Cookies.set('id', response.data.id, { expires: expiryDate });
            Cookies.set('username', response.data.user, { expires: expiryDate });
  
            // Preusmeri na domačo stran brez ponovnega nalaganja strani
            this.$router.push({ path: "/home" });
          } else {
            this.errorMessage = "Pri registraciji je prišlo do napake. Poskusite znova.";
          }
        } catch (error) {
          console.error("Registracijska napaka:", error);
  
          if (error.response && error.response.data && error.response.data.message) {
            this.errorMessage = error.response.data.message;
          } else {
            this.errorMessage = "Uporabniško ime je že v uporabi";
          }
        }
      },
    },
  };
  </script>
  
  <style scoped>
  @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap');
  
  html, body {
    margin: 0;
    padding: 0;
    height: 100%;
    width: 100%;
    overflow: hidden;
  }
  
  .register-page {
    position: relative;
    width: 100vw;
    height: 100vh;
    background: linear-gradient(135deg, #fbcfe8, #fde2e4, #fbd4e7, #fbe9e7);
    background-size: 400% 400%;
    animation: gradientBG 20s ease infinite;
    display: flex;
    justify-content: center;
    align-items: center;
    font-family: 'Poppins', sans-serif;
    color: #444;
    text-align: center;
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
    top: -50px;
    left: -50px;
  }
  
  .shape2 {
    width: 300px;
    height: 300px;
    background: #fde2e4;
    bottom: -100px;
    right: -100px;
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
    max-width: 500px;
    width: 90%;
    padding: 40px 30px;
    background: rgba(255, 255, 255, 0.25);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
    animation: fadeIn 2s ease forwards;
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
  
  .error-message {
    background-color: #f8d7da;
    color: #721c24;
    padding: 15px;
    border-radius: 10px;
    margin-bottom: 20px;
    font-weight: 400;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
  
  .register-form {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .form-group {
    width: 100%;
    margin-bottom: 20px;
    text-align: left;
  }
  
  .form-label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
    color: #333;
  }
  
  .form-input {
    width: 100%;
    padding: 12px 15px;
    border: 1px solid #ccc;
    border-radius: 10px;
    font-size: 1rem;
    transition: border-color 0.3s ease;
  }
  
  .form-input:focus {
    border-color: #8da4b9;
    outline: none;
    box-shadow: 0 0 5px rgba(138, 164, 185, 0.5);
  }
  
  .action-button {
    padding: 15px 40px;
    border-radius: 30px;
    text-decoration: none;
    font-weight: 600;
    font-size: 1.1rem;
    transition: background 0.3s ease, transform 0.3s ease;
    display: inline-block;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
    color: #fff;
    background-color: #8da4b9; /* Mehka pastelna modro-siva */
  }
  
  .action-button:hover {
    background-color: #7b93a9;
    transform: translateY(-3px);
  }
  
  .additional-links {
    margin-top: 20px;
    font-size: 1rem;
    color: #555;
  }
  
  .link {
    color: #f6a8c3;
    text-decoration: none;
    transition: color 0.3s ease;
  }
  
  .link:hover {
    color: #ee94b3;
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
    }
  
    .form-title {
      font-size: 1.75rem;
    }
  
    .action-button {
      width: 100%;
      text-align: center;
    }
  
    .buttons {
      flex-direction: column;
      gap: 15px;
    }
  }
  </style>
  