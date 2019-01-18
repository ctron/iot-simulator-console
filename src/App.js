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
                    <HorizontalNavMenuItem id="overview" title="Overview"/>
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
        };
        this.refreshData();
    }

    refreshData() {
        fetch(this.api + "/overview")
            .then(result => {
                return result.json()
            })
            .then(data => {
                console.log(data);
                this.setState({overview: data});
            })
    }

    componentDidMount() {
        this.interval = setInterval(() => this.refreshData(), 1000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    renderConsumers(tenant) {
        if (tenant.consumers == null) {
            return
        }
        return tenant.consumers.map(function (consumer, i) {
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

    renderProducers(tenant) {
        if (tenant.producers == null) {
            return
        }
        return tenant.producers.map(function (producer, i) {
            return (
                <ListViewItem
                    heading="producer"
                    description={producer.type + " / " + producer.protocol}
                    additionalInfo={[
                        <ListViewInfoItem>
                            <Icon type="fa" name="cubes"/>
                            <strong>{producer.replicas}</strong>
                        </ListViewInfoItem>,
                        <ListViewInfoItem>
                            <Icon type="fa" name="angle-double-up"/>
                            <strong>
                                {producer.messagesPerSecondConfigured} → {producer.messagesPerSecondScheduled ? producer.messagesPerSecondScheduled.toFixed(0) : "␀"} → {producer.messagesPerSecondSent ? producer.messagesPerSecondSent.toFixed(0) : "␀"}
                            </strong>
                        </ListViewInfoItem>
                    ]}
                >
                </ListViewItem>
            );
        })
    }

    render() {
        const o = this;
        return (
            <div>
                {
                    this.state.overview.tenants.map(function (tenant, i) {
                        return (
                            <ListGroup>
                                <h2>{tenant.name}</h2>
                                <ListView>
                                    {o.renderConsumers(tenant)}
                                    {o.renderProducers(tenant)}
                                </ListView>
                            </ListGroup>
                        )
                    })
                }
            </div>
        );
    }
}
