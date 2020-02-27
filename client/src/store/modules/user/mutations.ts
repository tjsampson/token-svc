import { MutationTree } from "vuex";
import { UserState } from "./types";

export const mutations: MutationTree<UserState> = {
  REGISTER_REQUEST(state, email) {
    state.loading = true;
  },
  REGISTER_SUCCESS(state, user) {
    state.loading = false;
    state.record = user;
  },
  REGISTER_FAILURE(state) {
    state.loading = false;
  },
  LOGIN_REQUEST(state, email) {
    state.loading = true;
  },
  LOGIN_SUCCESS(state, user) {
    state.loading = false;
    state.authenticated = true;
  },
  LOGIN_FAILURE(state) {
    state.loading = false;
  },
  LOGOUT_REQUEST(state) {
    state.loading = true;
  },
  LOGOUT_SUCCESS(state) {
    state.loading = false;
    state.authenticated = false;
  },
  LIST_USERS_REQUEST(state) {
    state.loading = true;
  },
  LIST_USERS_SUCCESS(state, users) {
    state.loading = false;
    state.users = users;
  },
  LIST_USERS_FAILURE(state) {
    state.loading = false;
    state.users = [];
  }
};
