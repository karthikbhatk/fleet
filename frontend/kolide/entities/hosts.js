import endpoints from "kolide/endpoints";

export default (client) => {
  return {
    destroy: (host) => {
      const { HOSTS } = endpoints;
      const endpoint = client._endpoint(`${HOSTS}/${host.id}`);

      return client.authenticatedDelete(endpoint);
    },

    refetch: (host) => {
      const { HOSTS } = endpoints;
      const endpoint = client._endpoint(`${HOSTS}/${host.id}/refetch`);

      return client
        .authenticatedPost(endpoint)
        .then((response) => response.host);
    },

    load: (hostID) => {
      const { HOSTS } = endpoints;
      const endpoint = client._endpoint(`${HOSTS}/${hostID}`);

      return client
        .authenticatedGet(endpoint)
        .then((response) => response.host);
    },

    loadAll: (
      page = 0,
      perPage = 100,
      selected = "",
      globalFilter = "",
      sortBy = []
    ) => {
      const { HOSTS, LABEL_HOSTS } = endpoints;

      // TODO: add this query param logic to client class
      const pagination = `page=${page}&per_page=${perPage}`;

      let orderKeyParam = "";
      let orderDirection = "";
      if (sortBy.length !== 0) {
        const sortItem = sortBy[0];
        orderKeyParam += `&order_key=${sortItem.id}`;
        orderDirection = `&order_direction=${sortItem.direction}`;
      }

      let searchQuery = "";
      if (globalFilter !== "") {
        searchQuery = `&query=${globalFilter}`;
      }

      let endpoint = "";
      const labelPrefix = "labels/";
      if (selected.startsWith(labelPrefix)) {
        const lid = selected.substr(labelPrefix.length);
        endpoint = `${LABEL_HOSTS(
          lid
        )}?${pagination}${searchQuery}${orderKeyParam}${orderDirection}`;
      } else {
        let selectedFilter = "";
        if (
          selected === "new" ||
          selected === "online" ||
          selected === "offline" ||
          selected === "mia"
        ) {
          selectedFilter = `&status=${selected}`;
        }
        endpoint = `${HOSTS}?${pagination}${selectedFilter}${searchQuery}${orderKeyParam}${orderDirection}`;
      }

      return client
        .authenticatedGet(client._endpoint(endpoint))
        .then((response) => {
          return response.hosts;
        });
    },
    transfer: () => {
      const { HOST_TRANSFER } = endpoints;
      const endpoint = client._endpoint(HOST_TRANSFER);
      return client
        .authenticatedGet(endpoint)
        .then((response: any) => response.secrets);
    },
  };
};
