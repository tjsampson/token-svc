<template>
  <div id="nav">
    <nav class="uk-navbar-container uk-margin" uk-navbar>
      <div class="uk-navbar-left">
        <div class="uk-navbar-item">
          <router-link to="/">Home</router-link>
        </div>
        <div class="uk-navbar-item">
          <router-link to="/register">Register</router-link>
        </div>
        <div class="uk-navbar-item">
          <router-link to="/users">Users</router-link>
        </div>
      </div>
      <div class="uk-navbar-right">
        <div class="uk-navbar-item">
          <button class="uk-button uk-button-default">
            <router-link to="/logout" v-if="authenticated == true"
              >Logout</router-link
            >
            <router-link to="/login" v-if="authenticated == false"
              >Login</router-link
            >
          </button>
        </div>
      </div>
    </nav>

    <alert></alert>
    <router-view />
  </div>
</template>

<script>
import Alert from "./Alert";
import { mapState } from "vuex";

export default {
  name: "topnav",
  computed: mapState({
    user: state => state.user,
    authenticated: state => state.user.authenticated
  }),

  components: {
    Alert
  }
};
</script>

<style lang="scss">
$global-link-color: #da7d02;

@import "uikit/src/scss/variables-theme.scss";
@import "uikit/src/scss/mixins-theme.scss";
@mixin hook-card() {
  color: #000;
}
@import "uikit/src/scss/uikit-theme.scss";

#nav {
  .uk-navbar-container {
    a {
      font-weight: bold;
      color: #2c3e50;

      &.router-link-exact-active {
        color: #42b983;
      }
    }
  }
}
</style>
