import Vue from "vue";
import VueRouter from "vue-router";
import { userService } from "@/_services";
// import jwtDecode from "jwt-decode";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "home" */ "../views/Home.vue"),
    meta: {
      public: true
    }
  },
  {
    path: "/login",
    name: "login",
    component: () => import(/* chunkName: "login" */ "../views/Login.vue"),
    meta: {
      public: true
    }
  },
  {
    path: "/logout",
    name: "logout",
    component: () => import(/* chunkName: "login" */ "../views/Logout.vue"),
    meta: {
      public: true
    }
  },
  {
    path: "/register",
    name: "register",
    component: () => import(/* chunkName: "users" */ "../views/Register.vue"),
    meta: {
      public: true
    }
  },
  {
    path: "/users",
    name: "users",
    component: () => import(/* chunkName: "users" */ "../views/Users.vue")
  },
  {
    path: "*",
    redirect: "/login"
  }
];

export const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes: routes
});

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.public)) {
    next();
  } else {
    // Private Routes (lets check some auth)
    let accessToken = localStorage.getItem(
      userService.userAccessTokenStorageKey
    );
    if (accessToken == null) {
      next({
        path: "/login",
        params: { nextUrl: to.fullPath }
      });
    } else {
      // let tokenData = jwtDecode(accessToken);
      next();
    }
  }
});

router.afterEach((to, from) => {
  if (to.matched.some(route => route.path == "/login")) {
    //console.log("TROY SAMPSON AFTER EACH TO LOGIN");
  }
  if (from.matched.some(route => route.path == "/login")) {
    //console.log("TROY SAMPSON AFTER EACH FROM LOGIN");
  }
});

export default router;
