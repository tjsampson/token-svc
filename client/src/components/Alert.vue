<template>
  <div v-if="alert.message" :class="`alert ${alert.type}`">
    <div uk-alert v-bind:class="[alert.type]">
      <a class="uk-alert-close" uk-close></a>
      {{ alert.message }}
    </div>
  </div>
</template>

<script>
export default {
  name: "alert",
  computed: {
    alert() {
      return this.$store.state.alert;
    }
  },
  watch: {
    $route(to, from) {
      // clear alert on location change
      this.$store.dispatch("alert/clear");
    }
  }
};
</script>
