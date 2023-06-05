<script>
    import GrafanaLogo from "./grafana-logo.svelte";

    let data = []
    let metric = ''
    let isLoading = false

    async function search() {
        if (metric === '') {
            data = []
            return
        }

        isLoading = true
        const res = await fetch(`/search?metric=${metric}`);
        data = await res.json();
        isLoading = false
    }

    function clearInput() {
        metric = ''
        data= []
    }
// todo: add support for url search params

</script>

<section class="section">
    <div class="container">
        <div class="level">
            <h1 class="level-item title has-text-centered">
                <GrafanaLogo/>Grafana metrics search
            </h1>
        </div>
        <form class="">
            <div class="level">
                <div class="level-item tile is-10 control has-icons-right">
                    <input class="input is-rounded is-large"
                           type="text"
                           placeholder="Type metric name here..."
                           bind:value={metric}
                    >
                    <span class="icon is-large is-right">
                        <button class="delete" on:click={clearInput}></button>
                    </span>
                </div>
                <div class="level-item tile is-2">
                    <button class="button is-primary is-rounded is-large {isLoading && 'is-loading'}" on:click={search}>Search</button>
                </div>
            </div>
        </form>
    </div>
</section>
{#if (data.length) !== 0}
    <section class="section">
        <table class="table is-striped is-fullwidth content is-medium">
            <thead>
            <tr>
                <th>Dashboard</th>
                <th>URL</th>
                <th><abbr title="Is used in grafana variables?">IsVar</abbr></th>
                <th>Panel</th>
            </tr>
            </thead>
            <tbody>
            {#each data as item}
                <tr>
                    <td class="is-clipped">{item.title}</td>
                    <td><a href="{item.url}" target="_blank">{item.url}</a></td>
                    <td>
                        {#if item.vars}
                            Yes
                        {/if}
                    </td>
                    <td>{item.panels}</td>
                </tr>
            {/each}
            </tbody>
        </table>
    </section>
{/if}
