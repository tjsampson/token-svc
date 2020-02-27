<template>
  <div class="uk-child-width-1-3@m uk-grid-small uk-grid-match" uk-grid>
    <div>
      <div
        class="uk-card uk-card-default uk-card-body uk-position-center uk-width-1-3 uk-margin"
      >
        <form @submit.prevent="submitRegister">
          <div class="uk-margin">
            <div class="uk-inline">
              <span class="uk-form-icon" uk-icon="icon: user"></span>
              <input
                class="uk-input"
                type="text"
                v-model="email"
                placeholder="email"
                name="email"
                :class="{ 'is-invalid': submitted && !email }"
              />
            </div>
          </div>

          <div class="uk-margin">
            <div class="uk-inline">
              <span class="uk-form-icon" uk-icon="icon: lock"></span>
              <input
                class="uk-input"
                type="password"
                v-model="password"
                placeholder="password"
                name="password"
                :class="{ 'is-invalid': submitted && !email }"
              />
            </div>
          </div>

          <div class="uk-margin">
            <div class="uk-inline">
              <span class="uk-form-icon" uk-icon="icon: lock"></span>
              <input
                class="uk-input"
                type="password"
                v-model="confirm_password"
                placeholder="confirm password"
                name="confirm_password"
                :class="{ 'is-invalid': submitted && !email }"
              />
            </div>
          </div>

          <button type="submit" class="uk-button uk-button-primary">
            Register
          </button>
          <div uk-spinner class="uk-hidden"></div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "registeruser",
  data() {
    return {
      email: "",
      password: "",
      confirm_password: "",
      submitted: false
    };
  },
  methods: {
    submitRegister(e) {
      this.submitted = true;
      const { email, password, confirm_password } = this;
      const { dispatch } = this.$store;
      if (email && password && confirm_password) {
        dispatch("user/register", { email, password, confirm_password });
      }
      this.submitted = false;
    }
  }
};
</script>

<style lang="scss" scoped>
.uk-card-default {
  background-color: #43596d;
}
</style>
