document.addEventListener('DOMContentLoaded', function () {
    const apiUrl = 'http://localhost:8080/api/v1/resources';
    let currentPage = 1;

    // Moved this event listener inside DOMContentLoaded
    document.getElementById('kindDropdown').addEventListener('change', function() {
        updateMetadataKeysDropdown(this.value);
    });

    // Setup filter form submission with proper event handling
    const filterForm = document.getElementById('filterForm');
    if (filterForm) {
        filterForm.removeEventListener('submit', handleFilterSubmit);
        filterForm.addEventListener('submit', handleFilterSubmit);
    }

    function handleFilterSubmit(event) {
        event.preventDefault();
        const filters = getCurrentFilters();
        fetchResources(1, filters);
    }

    function fetchResources(page, filters) {
        let url = `${apiUrl}?page=${page}`;

        for (const key in filters) {
            if (filters[key]) {
                url += `&${key}=${encodeURIComponent(filters[key])}`;
            }
        }

        console.log("Fetching URL:", url); // Add this line to check the constructed URL

        fetch(url)
            .then(response => response.json())
            .then(data => {
                displayResources(data.resources);
                setupPagination(data.page, data.has_more, filters);
            })
            .catch(error => console.error('Error fetching data:', error));
    }

    function displayResources(resources) {
        console.log("Displaying Resources:", resources); // Debugging line
        const tableBody = document.getElementById('resourceTable').getElementsByTagName('tbody')[0];
        tableBody.innerHTML = ''; // Clear existing rows

        resources.forEach(resource => {
            const row = tableBody.insertRow();
            row.innerHTML = `
            <td>${resource.kind}</td>
            <td>${resource.name}</td>
            <td>${resource.external_id}</td>
            <td>${resource.fetched_at}</td>
            <td>${resource.scanner}</td>
        `;
            row.addEventListener('click', () => {
                expandRow(row, resource);
            });
        });
    }


    function expandRow(row, resource) {
        const detailRow = row.nextElementSibling;
        if (detailRow && detailRow.classList.contains('detail-row')) {
            // Toggle existing detail row
            detailRow.remove();
        } else {
            // Insert new detail row
            row.insertAdjacentHTML('afterend', `
                <tr class="detail-row">
                    <td colspan="5">
                        <code><pre>${JSON.stringify(resource, null, 2)}</pre></code>
                    </td>
                </tr>
            `);
        }
    }

    function setupPagination(currentPage, hasMore, filters) {
        const pagination = document.getElementById('pagination');
        pagination.innerHTML = ''; // Clear existing buttons

        // Previous button
        if (currentPage > 1) {
            const prevBtn = createPaginationButton(currentPage - 1, 'Previous', filters);
            pagination.appendChild(prevBtn);
        }

        // Current page button
        const currentBtn = createPaginationButton(currentPage, currentPage.toString(), filters, true);
        pagination.appendChild(currentBtn);

        // Next button
        if (hasMore) {
            const nextBtn = createPaginationButton(currentPage + 1, 'Next', filters);
            pagination.appendChild(nextBtn);
        }
    }

    function createPaginationButton(page, text, filters, isActive = false) {
        const listItem = document.createElement('li');
        listItem.className = 'page-item';
        if (isActive) {
            listItem.classList.add('active');
        }

        const link = document.createElement('a');
        link.className = 'page-link';
        link.href = '#';
        link.textContent = text;

        link.addEventListener('click', function(event) {
            event.preventDefault();
            fetchResources(page, filters);
        });

        listItem.appendChild(link);
        return listItem;
    }


    function getCurrentFilters() {
        const kind = document.getElementById('kindDropdown') ? document.getElementById('kindDropdown').value : '';
        const metadataKey = document.getElementById('metadataKeys') ? document.getElementById('metadataKeys').value : '';
        const metadataValue = document.getElementById('metadataValue') ? document.getElementById('metadataValue').value : '';

        // Prepare the meta_data_eq filter if both key and value are provided
        let metaDataEq = '';
        if (metadataKey && metadataValue) {
            metaDataEq = `${metadataKey}=${metadataValue}`;
        }

        return {
            kind: kind,
            meta_data_eq: metaDataEq
        };
    }

    fetchResources(currentPage, {});
});

const kindMetadataMap = {
    "file_system": ["rootDirectory", "machineHost"],
    "aws_ec2": ["instance_id", "image_id", "private_dns_name", "instance_type", "architecture", "instance_lifecycle", "instance_state", "vpc_id"],
    "aws_ecr": ["repository_name", "arn", "registry_id", "repository_uri"],
    "aws_s3": ["bucket_name", "region", "arn"],
    "aws_rds": ["instance_id", "region", "arn"],
    "GitHubRepository": ["language", "stars",  "homepage", "organization", "company", "git_url", "owner_name", "owner_login"]
};

// Function to update metadata keys dropdown based on selected kind
function updateMetadataKeysDropdown(selectedKind) {
    const metadataKeysDropdown = document.getElementById('metadataKeys');
    metadataKeysDropdown.innerHTML = ''; // Clear existing options

    // Add an empty option as the default
    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.textContent = 'Select Metadata Key';
    metadataKeysDropdown.appendChild(defaultOption);

    if (kindMetadataMap[selectedKind]) {
        kindMetadataMap[selectedKind].forEach(key => {
            const option = document.createElement('option');
            option.value = key;
            option.textContent = key;
            metadataKeysDropdown.appendChild(option);
        });
    }
}
