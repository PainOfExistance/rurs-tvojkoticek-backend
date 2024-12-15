<template>
    <div class="login-container">
        <div v-if="errorMessage" class="alert alert-danger w-50 mx-auto" role="alert">{{ errorMessage }}</div>
        <div>
            <label for="username" class="form-label">Uporabniško ime:</label>
            <input
                type="text"
                id="username"
                class="form-control mb-2"
                required
                v-model="username"
                placeholder="Vnesi uporabniško ime..."
            >
        </div>
        <div>
            <label for="password" class="form-label">Geslo:</label>
            <input
                type="password"
                id="password"
                class="form-control mb-2"
                required
                v-model="password"
                placeholder="Vnesi geslo..."
            >
        </div>
        <button
            type="submit"
            class="btn btn-secondary mt-4 d-block mx-auto"
            @click="loginUser"
        >
            Prijava
        </button>
    </div>
</template>

<script>
import axios from "axios";
import Cookies from "js-cookie";
import router from "@/router";

export default {
    name: "Login",
    data() {
        return {
            username: "",
            password: "",
            errorMessage: "",
        };
    },
    methods: {
        async loginUser() {
            try {
                const response = await axios.post("http://localhost:8080/login", {
                    Username: this.username,
                    Password: this.password
                });
                const expiryDate = new Date();
                expiryDate.setDate(expiryDate.getDate() + 7);
                Cookies.set('id', response.data.id, { expires: expiryDate });
                Cookies.set('username', response.data.user, { expires: expiryDate });
                this.$router.push({ path: "/" }).then(() => {
                    window.location.reload();
                });
            } catch (error) {
                this.errorMessage = "Username or Password is incorrect";
            }
        },
    },
};
</script>

<style scoped>
.login-container {
  border: 1px solid #54627b;
  border-radius: 5px;
  padding: 30px 40px;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
}

.form-label {
  font-size: 17px;
  font-weight: bold;
}

.btn-secondary {
  background-color: #54627b;
}
</style>
