<template>
  <div class="dashboard-container">
</template>

<script>
export default {
  name: "Dashboard",
  components: { adminDashboard, editorDashboard },
  data() {
    return {};
  },
  computed: {
    ...mapGetters(["roles"]),
  },
  created() {
    if (!this.roles.includes("admin")) {
      this.currentRole = "editorDashboard";
    }
  },
};
</script>
