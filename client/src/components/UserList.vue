<template>
  <table class="uk-table uk-table-hover uk-table-divider">
    <thead>
      <tr>
        <th>ID</th>
        <th>UID</th>
        <th>Email</th>
      </tr>
    </thead>
    <tbody>
      <UserListItem
        v-for="item in user.users"
        :key="item.id"
        v-bind:user-details="item"
      />
    </tbody>
  </table>
</template>

<script>
import { mapState } from "vuex";
import UserListItem from "./UserListItem";

export default {
  name: "users",
  created() {
    this.$store.dispatch("user/list");
  },
  computed: mapState(["user"]),
  methods: {
    submitSearch(e) {}
  },
  components: {
    UserListItem
  }
};
</script>

<style scoped></style>
