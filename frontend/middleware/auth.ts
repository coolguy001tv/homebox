export default defineNuxtRouteMiddleware(async () => {
  const ctx = useAuthContext();
  const api = useUserApi();

  if (!ctx.isAuthorized()) {
    // In noAuth mode the status endpoint sets auth cookies and returns noAuth: true.
    // Check this before redirecting so a direct visit to /home works on first load.
    try {
      const publicApi = usePublicApi();
      const { data } = await publicApi.status();
      if (data?.noAuth) {
        ctx.enableNoAuth();
      }
    } catch {
      // status check failed — fall through to normal redirect logic
    }

    if (!ctx.isAuthorized()) {
      if (window.location.pathname !== "/") {
        console.debug("[middleware/auth] isAuthorized returned false, redirecting to /");
        return navigateTo("/");
      }
    }
  }

  if (!ctx.user) {
    console.log("Fetching user data");
    const { data, error } = await api.user.self();
    if (error) {
      if (window.location.pathname !== "/") {
        console.debug("[middleware/user] user is null and fetch failed, redirecting to /");
        return navigateTo("/");
      }
    }

    ctx.user = data.item;
  }
});
