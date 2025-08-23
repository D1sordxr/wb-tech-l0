import React from "react";
import './App.css'
import {useState} from "react";
import type {Order} from "./domain/Order.ts";
import {getMockOrder, getOrder} from "./api/Service.ts";

function App() {
    const mockOrderID = "b563feb7b2b84b6test";

    const [orderID, setOrderID] = useState(mockOrderID);
    const [orderData, setOrderData] = useState<Order | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);


    const fetchOrderData = async (orderId: string) => {
        try {
            if (orderId === mockOrderID) {
                const order = getMockOrder(orderId);
                setOrderData(order);
                return order;
            }
            
            setLoading(true);
            setError(null);
            const order: Order = await getOrder(orderId);
            setOrderData(order);
            return order;
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
            setError(errorMessage);
            console.error('Error fetching order:', error);
            throw error;
        } finally {
            setLoading(false);
        }
    };

    const handleSearch = async () => {
        if (!orderID.trim()) {
            setError('Please enter an order ID');
            return;
        }
        await fetchOrderData(orderID);
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setOrderID(e.target.value);
    };

    const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            handleSearch();
        }
    };

    return (
        <div className="container">
            <header>
                <h1>Order Search</h1>
                <p>Enter an order UID to retrieve its details</p>
            </header>

            <div className="search-form">
                <input
                    type="text"
                    value={orderID}
                    onChange={handleInputChange}
                    onKeyPress={handleKeyPress}
                    className="search-input"
                    placeholder="Enter order UID (e.g., b563feb7b2b84b6test)"
                    disabled={loading}
                />
                <button
                    onClick={handleSearch}
                    className="btn"
                    disabled={loading}
                >
                    {loading ? 'Searching...' : 'Search Order'}
                </button>
            </div>

            {error && (
                <div className="notification error">
                    {error}
                </div>
            )}

            {orderData && (
                <div className="order-data">
                    <div className="section">
                        <h2>Basic Information</h2>
                        <div className="info-grid">
                            <div className="info-item">
                                <div className="info-label">Order UID</div>
                                <div className="info-value">{orderData.order_uid}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Track Number</div>
                                <div className="info-value">{orderData.track_number}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Entry</div>
                                <div className="info-value">{orderData.entry}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Locale</div>
                                <div className="info-value">{orderData.locale}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Customer ID</div>
                                <div className="info-value">{orderData.customer_id}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Delivery Service</div>
                                <div className="info-value">{orderData.delivery_service}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Date Created</div>
                                <div className="info-value">{orderData.date_created}</div>
                            </div>
                        </div>
                    </div>

                    <div className="section">
                        <h2>Delivery Information</h2>
                        <div className="info-grid">
                            <div className="info-item">
                                <div className="info-label">Name</div>
                                <div className="info-value">{orderData.delivery.name}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Phone</div>
                                <div className="info-value">{orderData.delivery.phone}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Email</div>
                                <div className="info-value">{orderData.delivery.email}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Address</div>
                                <div className="info-value">{orderData.delivery.address}, {orderData.delivery.city}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Zip Code</div>
                                <div className="info-value">{orderData.delivery.zip}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Region</div>
                                <div className="info-value">{orderData.delivery.region}</div>
                            </div>
                        </div>
                    </div>

                    <div className="section">
                        <h2>Payment Information</h2>
                        <div className="info-grid">
                            <div className="info-item">
                                <div className="info-label">Transaction</div>
                                <div className="info-value">{orderData.payment.transaction}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Currency</div>
                                <div className="info-value">{orderData.payment.currency}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Provider</div>
                                <div className="info-value">{orderData.payment.provider}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Amount</div>
                                <div className="info-value">${orderData.payment.amount.toFixed(2)}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Payment Date</div>
                                <div className="info-value">{orderData.payment.payment_dt}</div>
                            </div>
                            <div className="info-item">
                                <div className="info-label">Bank</div>
                                <div className="info-value">{orderData.payment.bank}</div>
                            </div>
                        </div>
                    </div>

                    <div className="section">
                        <h2>Order Items</h2>
                        <table className="items-table">
                            <thead>
                            <tr>
                                <th>Product</th>
                                <th>Brand</th>
                                <th>Price</th>
                                <th>Sale</th>
                                <th>Total</th>
                                <th>Status</th>
                            </tr>
                            </thead>
                            <tbody>
                            {orderData.items.map((item, index) => (
                                <tr key={index}>
                                    <td>{item.name}</td>
                                    <td>{item.brand}</td>
                                    <td>${item.price.toFixed(2)}</td>
                                    <td>{item.sale}%</td>
                                    <td>${item.total_price.toFixed(2)}</td>
                                    <td>
                                        <span className={`status-badge ${item.status === 202 ? 'status-completed' : 'status-pending'}`}>
                                            {item.status === 202 ? 'Completed' : 'Pending'}
                                        </span>
                                    </td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}
        </div>
    )
}

export default App