<template>
  <div class="d-flex justify-content-center mt-4">
    <!-- Profile Information Table -->
    <div class="profile-info-container p-4 border rounded shadow-sm bg-light me-3" style="height: fit-content">
      <h2 class="text-center mb-4">Profile</h2>
      <div class="mb-4">
        <h5>Username:</h5>
        <p>{{ username }}</p>

        <h5>Email:</h5>
        <p>{{ email }}</p>
      </div>
    </div>

    <!-- Change Password Table -->
    <div class="password-change-container p-4 border rounded shadow-sm bg-light" style="height: fit-content">
      <h4 class="text-center mb-4">Change Password</h4>
      <form @submit.prevent="handleChangePassword">
        <div class="form-group mb-3">
          <label class="form-label fw-bold" style="font-size: 0.9rem;">New Password:</label>
          <input type="password" v-model="newPassword" class="form-control" required>
        </div>

        <div class="form-group mb-3">
          <label class="form-label fw-bold" style="font-size: 0.9rem;">Confirm New Password:</label>
          <input type="password" v-model="confirmPassword" class="form-control" required>
        </div>

        <button type="submit" class="btn btn-primary mt-3">Change Password</button>
      </form>

      <p v-if="message" class="alert alert-info mt-3">{{ message }}</p>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Cookies from 'js-cookie';
import router from "@/router";

export default {
  data() {
    return {
      isLoggedIn: false,
      username: '',
      email: '',
      newPassword: '',
      confirmPassword: '',
      message: ''
    };
  },
  mounted() {
    this.checkLoginStatus();
    this.getUserInfo();
  },
  methods: {
    getUserInfo() {
      if (!this.isLoggedIn) {
        router.push({ path: '/login' })
      }
      axios.get('http://localhost:8080/profile', {
        params: {
          username: Cookies.get('username')
        }
      })
          .then(response => {
            console.log('User info:', response.data);
            this.username = response.data.username;
            this.email = response.data.email;
          })
          .catch(error => {
            console.error('Error fetching user info:', error);
            this.message = error.response?.data?.message || 'Failed to load user data.';
          });
    },
    handleChangePassword() {
      if (this.newPassword !== this.confirmPassword) {
        this.message = "New passwords do not match.";
        return;
      }

      axios.post('http://localhost:8080/changePassword', {
        username: this.username,
        newPassword: this.newPassword
      })
          .then(response => {
            this.message = response.data.message || 'Password changed successfully!';
            this.clearForm();
          })
          .catch(error => {
            this.message = error.response?.data?.message || 'An error occurred while changing the password.';
            console.error('Error changing password:', error);
          });
    },
    clearForm() {
      this.newPassword = '';
      this.confirmPassword = '';
    },
    checkLoginStatus() {
      this.isLoggedIn = !!Cookies.get('id');
    }
  }
};
</script>

<style scoped>
.profile-info-container {
  width: 200px;
}

.profile-info-container,
.password-change-container {
  background-color: #f8f9fa;
  border-radius: 8px;
}
</style>
