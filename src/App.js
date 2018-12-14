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
            <ListView>
                {
                    this.state.overview.tenants.map(function (tenant, i) {
                        return (
                            <ListViewItem
                                heading={tenant.name}
                                additionalInfo={[
                                    <ListViewInfoItem>
                                        <Icon type="fa" name="angle-double-up"/>
                                        <strong>{tenant.statistics.telemetry.producer.messagesPerSecond} / {tenant.statistics.event.producer.messagesPerSecond}</strong>
                                        Producer
                                    </ListViewInfoItem>,
                                    <ListViewInfoItem>
                                        <Icon type="fa" name="angle-double-down"/>
                                        <strong>{tenant.statistics.telemetry.consumer.replicas} / {tenant.statistics.event.consumer.replicas}</strong>
                                        Consumer
                                    </ListViewInfoItem>,
                                ]}
                            >
                            </ListViewItem>
                        );
                    })
                }
            </ListView>
        )
    }
}
