import authUser from "@/components/authUser.vue";
import pastebinForm from "../components/pastebinForm.vue";
import pastebinList from "@/components/pastebinList.vue";
import pastebinDetail from "@/components/pastebinDetail.vue";
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/auth",
      name: "authUser",
      component: authUser,
    },
    {
      path: "/newPast",
      name: "createPast",
      component: pastebinForm
    },
    {
      path: "/past/:id",
      name: "pastDetails",
      component: pastebinDetail
    },
    {
      path: "/pasteList",
      name: "pasteList",
      component: pastebinList
    }
  ],
});

console.log(router);

export default router;