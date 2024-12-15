<template>
    <div class="register-container">
        <div v-if="errorMessage" class="alert alert-danger" role="alert">{{ errorMessage }}</div>
        <div>
            <label for="username" class="form-label">Uporabniško ime:</label>
            <input
                type="text"
                id="username"
                class="form-control mb-3"
                required
                v-model="username"
                placeholder="Vnesi uporabniško ime..."
            >
        </div>
        <div>
            <label for="email" class="form-label">Elektronski naslov:</label>
            <input
                type="email"
                id="email"
                class="form-control mb-2"
                required
                v-model="email"
                placeholder="Vnesi elektronski naslov..."
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
        <div>
            <label for="repeat_password" class="form-label">Ponovi geslo:</label>
            <input
                type="password"
                id="repeat_password"
                class="form-control mb-2"
                required
                v-model="repeat_password"
                placeholder="Ponovno vnesi geslo..."
            >
        </div>
        <button class="btn btn-secondary d-block mx-auto mt-4" @click="addUser">Registracija</button>
    </div>
</template>

<script>
import axios from "axios";

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
            if (this.password === this.repeat_password) {
                if (this.username && this.password && this.email) {
                    try {
                        const response = await axios.post("http://localhost:8080/register", {
                            Username: this.username,
                            Email: this.email,
                            Password: this.password,
                        });
                        this.$router.push({ path: "/login" }).then(() => {
                            window.location.reload();
                        });
                    } catch (error) {
                        this.errorMessage = "Username is already in use";
                    }
                } else {
                    this.errorMessage = "Please fill in all fields";
                }
            } else {
                this.errorMessage = "Passwords do not match";
            }
        }
    }
};
</script>

<style scoped>
.register-container {
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
