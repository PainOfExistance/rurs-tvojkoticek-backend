<template>
  <div>
    <h1 class="text-center">Vprašaj!</h1>

    <div v-if="isLoggedIn" class="post-add-container">
      <form @submit.prevent="submitAQuestion">
        <label class="form-label">Avtor:</label>
        <input id="authorInput" class="form-control mb-2" placeholder="Vnesite vaše ime..." required style="font-size: 0.9rem;" v-model="formData.name">

        <label class="form-label">Vaše vprašanje:</label>
        <textarea id="storyInput" type="text" class="post-input-area form-control mb-2" rows="5" placeholder="Opišite vaš problem in kaj vas zanima..." required style="font-size: 0.9rem;" v-model="formData.qustion"></textarea>

        <button type="submit" class="btn btn-danger mt-4 d-block mx-auto">Vprašaj</button>
      </form>
    </div>

    <div v-else class="text-center">
      <p class="mt-4 not-logged-in-info">V kolikor želite postaviti vprašanje morate biti v sistem prijavljeni.</p>
      <div class="row row-cols-2 mt-4">
        <div class="col text-center">
          <button class="btn btn-secondary"><router-link class="button-link-text" to="/login">Prijava</router-link></button>
        </div>
        <div class="col text-center">
          <button class="btn btn-secondary"><router-link class="button-link-text" to="/register">Registracija</router-link></button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Cookies from 'js-cookie';
import axios from "axios";

export default {
  data() {
    return {
      isLoggedIn: false,
      formData: {
        name: '',
        qustion: '',
      },
    };
  },
  methods: {
    checkLoginStatus() {
      this.isLoggedIn = !!Cookies.get('id');
    },
    submitAQuestion() {
      console.log(this.formData.name);
      axios.post('http://localhost:8080/post', {
        username: this.formData.name,
        problem: this.formData.qustion,
        date: new Date().toISOString().split('T')[0],
      }).then(response => {
        if (response.status === 200) {
          this.$router.push({ name: 'Posts' });
        }
      }).catch(error => console.error(error));
    }
  },
  mounted() {
    this.checkLoginStatus();
  },
};
</script>

<style>
.post-add-container {
  border: 1px solid #54627b;
  border-radius: 5px;
  padding: 30px 40px;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
}

.post-input-area {
  max-height: 40%;
}

.not-logged-in-info {
  font-size: 22px;
  font-weight: bold;
}

.form-label {
  font-size: 17px;
  font-weight: bold;
}

.button-link-text {
  color: white;
  text-decoration: none;
}

.btn-danger {
  background-color: #fbb2a3 !important;
  border-color: #fbb2a3 !important;
}

.btn-danger:hover {
  background-color: #fd9c8c !important;
  border-color: #fd9c8c !important;
}

.btn-secondary {
  background-color: #54627b;
}
</style>
