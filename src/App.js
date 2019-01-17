import $ from 'jquery';
import React from "react";
import {
    AggregateStatusCount,
    Card,
    CardTitle,
    CardGrid,
    CardBody,
    Icon,
    Row,
    Col,
    AggregateStatusNotifications,
    AggregateStatusNotification,
    Masthead,
    MastheadCollapse,
    MastheadDropdown,
    MenuItem,
    ListGroup,
    ListView,
    ListViewItem, ListViewInfoItem,
} from "patternfly-react";

import "./App.css"

import HorizontalNav from "patternfly-react/dist/esm/components/HorizontalNav/HorizontalNav";
import HorizontalNavMenu from "patternfly-react/dist/esm/components/HorizontalNav/HorizontalNavMenu";
import HorizontalNavMenuItem from "patternfly-react/dist/esm/components/HorizontalNav/HorizontalNavMenuItem";

class App extends React.Component {
    render() {
        return <div>
            <HorizontalNav>
                <Masthead
                    title="IoT simulator console"
                    navToggle={false}
                    thin
                />
                <HorizontalNavMenu>
                    <HorizontalNavMenuItem id="overview" title="Overview"></HorizontalNavMenuItem>
                </HorizontalNavMenu>
            </HorizontalNav>
            <div className="container-fluid">
                <Home/>
            </div>
        </div>
    }
}

export default App;

class Home extends React.Component {

    constructor(props) {
        super(props);
        this.api = "//" + window.location.host + "/api";
        this.state = {
            overview: {
                tenants: []
            },
        }
        this.refreshData();
    }

    refreshData() {
        fetch(this.api + "/overview")
            .then(result => {
                return result.json()
            })
            .then(data => {
                console.log(data)
                this.setState({overview: data})
            })
    }

    componentDidMount() {
        this.interval = setInterval(() => this.refreshData(), 1000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    render() {
        return (
            <div>
                {
                    this.state.overview.tenants.map(function (tenant, i) {
                        return (
                            <ListGroup>
                                <h2>{tenant.name}</h2>
                                <ListView>
                                    {
                                        tenant.consumers.map(function (consumer, i) {
                                            return (
                                                <ListViewItem
                                                    heading="consumer"
                                                    description={consumer.type}
                                                    additionalInfo={[
                                                        <ListViewInfoItem>
                                                            <Icon type="fa" name="cubes"/>
                                                            <strong>{consumer.replicas}</strong>
                                                        </ListViewInfoItem>,
                                                        <ListViewInfoItem>
                                                            <Icon type="fa" name="angle-double-down"/>
                                                            <strong>{consumer.messagesPerSecond}</strong>
                                                        </ListViewInfoItem>
                                                    ]}
                                                >
                                                </ListViewItem>
                                            )
                                        })
                                    }
                                </ListView>
                            </ListGroup>
                        )
                    })
                }
            </div>
        )
    }
}
