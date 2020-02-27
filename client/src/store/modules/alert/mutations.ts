import { MutationTree } from "vuex";
import { AlertState } from "./types";

export const mutations: MutationTree<AlertState> = {
  success(state, message) {
    state.type = "uk-alert-success";
    state.message = message;
  },
  warning(state, message) {
    state.type = "uk-alert-warning";
    state.message = message;
  },
  error(state, message) {
    state.type = "uk-alert-danger";
    state.message = message;
  },
  clear(state) {
    state.type = null;
    state.message = null;
  }
};
