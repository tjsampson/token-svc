import Vue from "vue";
import Vuex, { StoreOptions } from "vuex";
import { RootState } from "@/store/root";
import user from "@/store/modules/user";
import alert from "@/store/modules/alert";

Vue.use(Vuex);

const store: StoreOptions<RootState> = {
  modules: {
    user: user,
    alert: alert
  }
};

export default new Vuex.Store<RootState>(store);
